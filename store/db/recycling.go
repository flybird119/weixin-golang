package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//获取云店回收信息
func AccessStoreRecylingInfo(recyling *pb.Recyling) error {
	//首先获取店铺回收信息
	query := "select id,store_id,appoint_times,status,summary,qrcode_url,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint from recyling where store_id='%s'"
	query = fmt.Sprintf(query, recyling.StoreId)
	log.Debug(query)
	var times []*pb.RecylingAppointTime
	var appoint_times string
	err := DB.QueryRow(query).Scan(&recyling.Id, &recyling.StoreId, &appoint_times, &recyling.Status, &recyling.Summary, &recyling.QrcodeUrl, &recyling.CreateAt, &recyling.UpdateAt)
	//如果没有找到店铺数据，那么初始化店铺回收数据
	if err == sql.ErrNoRows {
		query = "insert into recyling (store_id) values('%s') returning id,store_id,appoint_times,status,summary,qrcode_url,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint"
		query = fmt.Sprintf(query, recyling.StoreId)
		err = DB.QueryRow(query).Scan(&recyling.Id, &recyling.StoreId, &appoint_times, &recyling.Status, &recyling.Summary, &recyling.QrcodeUrl, &recyling.CreateAt, &recyling.UpdateAt)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}

	if err := json.Unmarshal([]byte(appoint_times), &times); err != nil {
		log.Debug(err)

	}

	recyling.AppointTimes = times
	return nil
}

//提交预约订单接口
func UserSubmitRecylingOrder(recylingOrder *pb.RecylingOrder) error {

	searchOrder := &pb.RecylingOrder{UserId: recylingOrder.UserId}
	err := UserAccessPendingRecylingOrder(searchOrder)
	if err != nil {
		log.Error(err)
		return err
	}
	if searchOrder.Id != "" {
		return errors.New("alreadyExists")
	}

	query := "insert into recyling_order (store_id,school_id,lp_user_id,images,remark,addr,mobile,appoint_start_at,appoint_end_at) values('%s','%s','%s','%s','%s','%s','%s',to_timestamp(%d),to_timestamp(%d)) returning id"

	//转化图片
	imagesBytes, err := json.Marshal(recylingOrder.Images)
	var imagesStr string
	if err != nil || recylingOrder.Images == nil {
		imagesStr = "[]"
	} else {
		imagesStr = string(imagesBytes)
	}
	query = fmt.Sprintf(query, recylingOrder.StoreId, recylingOrder.SchoolId, recylingOrder.UserId, imagesStr, recylingOrder.Remark, recylingOrder.Addr, recylingOrder.Mobile, recylingOrder.AppointStartAt, recylingOrder.AppointEndAt)

	log.Debug(query)

	err = DB.QueryRow(query).Scan(&recylingOrder.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//查看预约中的回收订单接口
func UserAccessPendingRecylingOrder(recylingOrder *pb.RecylingOrder) error {

	query := "select id,store_id,school_id,lp_user_id,images,state,remark,addr,mobile,extract(epoch from appoint_start_at)::bigint,extract(epoch from appoint_end_at)::bigint,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint from recyling_order where state in(1,2) and lp_user_id='%s'"
	log.Debug(query)
	query = fmt.Sprintf(query, recylingOrder.UserId)
	var images []*pb.RecylingImage
	var imageStr string
	err := DB.QueryRow(query).Scan(&recylingOrder.Id, &recylingOrder.StoreId, &recylingOrder.SchoolId, &recylingOrder.UserId, &imageStr, &recylingOrder.State, &recylingOrder.Remark, &recylingOrder.Addr, &recylingOrder.Mobile, &recylingOrder.AppointStartAt, &recylingOrder.AppointEndAt, &recylingOrder.CreateAt, &recylingOrder.UpdateAt)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		log.Error(err)
		return err
	}
	if err := json.Unmarshal([]byte(imageStr), &images); err != nil {
		log.Debug(err)
	}
	recylingOrder.Images = images

	return nil
}

//设置云店回收信息
func UpdateStoreRecylingInfo(recyling *pb.Recyling) error {
	query := "update recyling set update_at=now()"

	var condition string
	if recyling.Status != 0 {
		condition += fmt.Sprintf(",status=%d", recyling.Status)
	}
	if recyling.QrcodeUrl != "" {
		condition += fmt.Sprintf(",qrcode_url='%s'", recyling.QrcodeUrl)
	}
	if recyling.Summary != "" {
		condition += fmt.Sprintf(",summary='%s'", recyling.Summary)
	}
	if recyling.AppointTimes != nil {
		//转化图片
		timeBytes, err := json.Marshal(recyling.AppointTimes)
		var timeStr string
		if err != nil {
			timeStr = "[]"
		} else {
			timeStr = string(timeBytes)
		}
		condition += fmt.Sprintf(",appoint_times='%s'", timeStr)
	}
	condition += fmt.Sprintf(" where id='%s'", recyling.Id)
	query += condition
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//获取云店回收订单列表
func GetStoreRecylingOrderList(recylingOrder *pb.RecylingOrder) (models []*pb.RecylingOrder, err error, totalCount int64) {
	query := "select id,store_id,school_id,lp_user_id,images,state,remark,addr,mobile,extract(epoch from appoint_start_at)::bigint,extract(epoch from appoint_end_at)::bigint,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint from recyling_order where 1=1"
	queryCount := "select count(*) from recyling_order where 1=1 "
	var condition string
	if recylingOrder.State != 0 {
		condition += fmt.Sprintf(" and state=%d", recylingOrder.State)
	}
	if recylingOrder.SortBy == "" {
		recylingOrder.SortBy = " appoint_start_at"
	}
	if recylingOrder.SequenceBy == "" {
		recylingOrder.SequenceBy = " desc"
	}
	condition += fmt.Sprintf(" and school_id='%s' and store_id='%s'", recylingOrder.SchoolId, recylingOrder.StoreId)

	if recylingOrder.Page <= 0 {
		recylingOrder.Page = 1
	}
	if recylingOrder.Size <= 0 {
		recylingOrder.Size = 15
	}
	//查询符合条件数据的长度
	queryCount += condition
	log.Debug(queryCount)
	err = DB.QueryRow(queryCount).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}
	//如果数据长度为零，直接返回 ，下面的查询语句就不执行了
	if totalCount == 0 {

		return
	}
	condition += fmt.Sprintf("  order by %s %s OFFSET %d LIMIT %d ", recylingOrder.SortBy, recylingOrder.SequenceBy, (recylingOrder.Page-1)*recylingOrder.Size, recylingOrder.Size)
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return models, nil, 0
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	//遍历rows
	for rows.Next() {
		order := &pb.RecylingOrder{}
		models = append(models, order)
		var images []*pb.RecylingImage
		var imageStr string
		err = rows.Scan(&order.Id, &order.StoreId, &order.SchoolId, &order.UserId, &imageStr, &order.State, &order.Remark, &order.Addr, &order.Mobile, &order.AppointStartAt, &order.AppointEndAt, &order.CreateAt, &order.UpdateAt)
		if err != nil {
			log.Error(err)
			return
		}
		if err := json.Unmarshal([]byte(imageStr), &images); err != nil {
			log.Debug(err)
		}
		order.Images = images

	}
	return

}

//更改回收订单
func UpdateRecylingOrder(recylingOrder *pb.RecylingOrder) error {
	//更改回收状态
	//回收状态 1 待处理 2 搁置中 3 已完成
	query := "update recyling_order set update_at=now()"

	if recylingOrder.State == 2 || recylingOrder.State == 3 {
		query += fmt.Sprintf(",state=%d", recylingOrder.State)
	}

	if recylingOrder.SellerRemark != "" {
		query += fmt.Sprintf(",seller_remark='%s'", recylingOrder.SellerRemark)
	}

	if recylingOrder.AppointStartAt != 0 && recylingOrder.AppointEndAt != 0 {
		query += fmt.Sprintf(",appoint_start_at=to_timestamp(%d),appoint_end_at=to_timestamp(%d)", recylingOrder.AppointStartAt, recylingOrder.AppointEndAt)
	}

	query += fmt.Sprintf(" where id='%s' and store_id='%s'", recylingOrder.Id, recylingOrder.StoreId)
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
