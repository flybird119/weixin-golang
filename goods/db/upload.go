package db

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//增加批量上传记录
func AddBatchUpload(model *pb.GoodsBatchUploadModel) error {
	query := "insert into goods_batch_upload (store_id,type,discount,storehouse_id,shelf_id,floor_id,origin_file,origin_filename) values('%s',%d,%d,'%s','%s','%s','%s','%s') returning id,extract(epoch from create_at)::bigint"
	query = fmt.Sprintf(query, model.StoreId, model.Type, model.Discount, model.StorehouseId, model.ShelfId, model.FloorId, model.OriginFile, model.OriginFilename)
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&model.Id, &model.CreateAt)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//更改批量上传记录
func UpdateBatchUpload(model *pb.GoodsBatchUploadModel) error {
	query := "update goods_batch_upload set update_at=now()"
	var condition string
	if model.SuccessNum != 0 {
		condition += fmt.Sprintf(",success_num=%d", model.SuccessNum)
	}
	if model.FailedNum != 0 {
		condition += fmt.Sprintf(",failed_num=%d", model.FailedNum)
	}
	if model.State != 0 {
		condition += fmt.Sprintf(",state=%d", model.State)
	}
	if model.ErrorReason != "" {
		condition += fmt.Sprintf(",error_reason='%s'", model.ErrorReason)
	}
	if model.ErrorFile != "" {
		condition += fmt.Sprintf(",error_file='%s'", model.ErrorFile)
	}
	if model.CompleteAt != 0 {
		condition += fmt.Sprintf(",complete_at=now()")
	}
	condition += fmt.Sprintf(" where id='%s'", model.Id)
	query += condition
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//获取批量上传数据
func GoodsBactchUploadList(model *pb.GoodsBatchUploadModel) (models []*pb.GoodsBatchUploadModel, err error, totalCount int64) {
	query := "select %s from goods_batch_upload where 1=1 and store_id='%s' %s"

	param := "id,store_id,success_num,failed_num,state,type,discount,storehouse_id,shelf_id,floor_id,error_reason,origin_file,origin_filename,error_file,extract(epoch from create_at)::bigint,extract(epoch from complete_at)::bigint,extract(epoch from update_at)::bigint"
	paramCount := "count(*)"

	var idParam string
	if model.Id != "" {
		idParam += fmt.Sprintf(" and id='%s'", model.Id)
	}

	now := time.Now()
	now = now.AddDate(0, 0, -30)

	idParam += fmt.Sprintf(" and extract(epoch from create_at)::bigint>%d", now.Unix())
	queryCount := fmt.Sprintf(query, paramCount, model.StoreId, idParam)
	query = fmt.Sprintf(query, param, model.StoreId, idParam)

	//查询条数
	log.Debug(queryCount)
	err = DB.QueryRow(queryCount).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}
	if totalCount <= 0 {
		return
	}
	//查询附加条件
	var condition string

	if model.Page <= 0 {
		model.Page = 1
	}

	if model.Size <= 0 {
		model.Size = 15
	}
	condition += fmt.Sprintf(" OFFSET %d LIMIT %d ", (model.Page-1)*model.Size, model.Size)
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	//遍历数据
	for rows.Next() {
		result := &pb.GoodsBatchUploadModel{}
		models = append(models, result)
		var completeAt sql.NullInt64
		//id,store_id,success_num,failed_num,state,type,discount,storehouse_id,shelf_id,floor_id,error_reason,origin_file,origin_filename,error_file,extract(epoch from create_at)::bigint,extract(epoch from complete_at)::bigint,extract(epoch from update_at)::bigint
		err = rows.Scan(&result.Id, &result.StoreId, &result.SuccessNum, &result.FailedNum, &result.State, &result.Type, &result.Discount, &result.StorehouseId, &result.ShelfId, &result.FloorId, &result.ErrorReason, &result.OriginFile, &result.OriginFilename, &result.ErrorFile, &result.CreateAt, &completeAt, &result.UpdateAt)
		if err != nil {
			log.Error(err)
			return
		}
		if completeAt.Valid {
			result.CompleteAt = completeAt.Int64
		}
		//获取位置信息
		err = getLocationInfo(result)
		if err != nil {
			return
		}
	}
	return
}

//获取位置信息
func getLocationInfo(model *pb.GoodsBatchUploadModel) error {
	var id, name string
	query := fmt.Sprintf("select id,name from location where id in ('%s','%s','%s')", model.StorehouseId, model.ShelfId, model.FloorId)
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil && err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Error(err)
			return err
		}
		if id == model.StorehouseId {
			model.StorehouseName = name
		} else if id == model.ShelfId {
			model.ShelfName = name
		} else {
			model.FloorName = name
		}
	}
	return nil
}
