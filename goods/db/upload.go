package db

import (
	"fmt"

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
