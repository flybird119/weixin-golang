package db

import (
	"database/sql"
	"strings"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//查找数据库过期的订单
func FindAllExpireOrder() (orders []*pb.Order, err error) {
	//找出半小时未付款的订单
	nowTime := time.Now().Unix()
	//订单截止时间
	expireTimestamp := nowTime - 30*60
	query := "select id,order_status,user_id from orders where order_status=0 and extract(epoch from order_at)::integer < $1"
	log.Debugf("select id,order_status,user_id from orders where order_status=0 and extract(epoch from order_at)::integer < %d", expireTimestamp)
	rows, err := DB.Query(query, expireTimestamp)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		order := &pb.Order{}
		orders = append(orders, order)
		err = rows.Scan(&order.Id, &order.OrderStatus, &order.UserId)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

//统计线上交易额
func OnlineGoodsSalesStatistic(goodsSalesStatisticModel *pb.GoodsSalesStatisticModel) error {
	//统计渠道销售额
	//1.0订单量 和 相对应的销售额
	//查询付过款 ，并且在一定时间阶段的学校订单
	//1.1 统计渠道销售额
	query := "select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from order_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2 and (order_status >0 and order_status<>8)  group by pay_channel"
	log.Debugf("select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from order_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s' and (order_status >0 and order_status<>8)  group by pay_channel", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	rows, err := DB.Query(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	defer rows.Close()
	var alipay_order_num, alipay_order_fee, wechat_order_num, wechat_order_fee, fee int64
	var pay_channel string
	for rows.Next() {
		err = rows.Scan(&fee, &pay_channel)
		if err != nil {
			log.Error(err)
			return err
		}
		if strings.Contains(pay_channel, "alipay") {
			alipay_order_num += 1
			alipay_order_fee += fee
		} else {
			wechat_order_num += 1
			wechat_order_fee += fee
		}
	}
	goodsSalesStatisticModel.AlipayOrderNum = alipay_order_num
	goodsSalesStatisticModel.AlipayOrderFee = alipay_order_fee
	goodsSalesStatisticModel.WechatOrderNum = wechat_order_num
	goodsSalesStatisticModel.WechatOrderFee = wechat_order_fee

	//1.2统计线上新书旧书销售额
	query = "select sum(oi.price),oi.type from orders_item oi join orders o on oi.orders_id=o.id and to_char(to_timestamp(extract(epoch from o.pay_at )::integer), 'YYYY-MM-DD')=$1 and o.school_id=$2  group by oi.type"
	log.Debugf("select sum(oi.price),oi.type from orders_item oi join orders o on oi.orders_id=o.id and to_char(to_timestamp(extract(epoch from o.pay_at )::integer), 'YYYY-MM-DD')='%s' and o.school_id='%s'  group by oi.type", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	rows, err = DB.Query(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	defer rows.Close()
	var oldbookSales, newbookSales, bookType int64
	var price sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&price, &bookType)
		if err != nil {
			log.Error(err)
			return err
		}
		if bookType == 0 {
			if price.Valid {
				newbookSales += price.Int64
			} else {
				newbookSales += 0
			}

		} else {
			if price.Valid {
				oldbookSales += price.Int64
			} else {
				oldbookSales += 0
			}
		}

	}
	goodsSalesStatisticModel.OnlineNewBookSalesFee = newbookSales
	goodsSalesStatisticModel.OnlineOldBookSalesFee = oldbookSales

	//2.0 统计日发送订单
	//*** 2.1
	query = "select count(*) from orders where to_char(to_timestamp(extract(epoch from deliver_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2"
	log.Debugf("select count(*) from orders where to_char(to_timestamp(extract(epoch from deliver_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s'", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	err = DB.QueryRow(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId).Scan(&goodsSalesStatisticModel.SendOrderNum)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	//1.2统计线上新书旧书销售额

	//3.0 统计日申请售后数量
	query = "select count(*) from orders where to_char(to_timestamp(extract(epoch from after_sale_apply_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2"
	log.Debugf("select count(*) from orders where to_char(to_timestamp(extract(epoch from after_sale_apply_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s'", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	err = DB.QueryRow(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId).Scan(&goodsSalesStatisticModel.AfterSaleNum)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	//4.0 统计日处理售后数量和费用
	query = "select count(*), sum(total_fee) from orders where to_char(to_timestamp(extract(epoch from after_sale_apply_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2"
	log.Debugf("select count(*),sum(total_fee) from orders where to_char(to_timestamp(extract(epoch from after_sale_apply_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s'", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	var totalFee sql.NullInt64
	err = DB.QueryRow(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId).Scan(&goodsSalesStatisticModel.AfterSaleHandledNum, &totalFee)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	if totalFee.Valid {
		goodsSalesStatisticModel.AfterSaleHandledFee = totalFee.Int64
	}
	return nil
}

//线下销售统计
func OfflineGoodsSalesStatistic(goodsSalesStatisticModel *pb.GoodsSalesStatisticModel) error {
	//线下销售数据统计
	//统计订单量
	query := "select count(*) from retail where to_char(to_timestamp(extract(epoch from create_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2"
	log.Debugf("select count(*) from retail where to_char(to_timestamp(extract(epoch from create_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s'", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	err := DB.QueryRow(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId).Scan(&goodsSalesStatisticModel.OfflineOrderNum)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	//统计新书旧书销售额
	query = "select sum(ri.price),ri.type from retail_item ri join retail r on ri.retail_id=r.id and to_char(to_timestamp(extract(epoch from r.create_at )::integer), 'YYYY-MM-DD')=$1 and r.school_id=$2  group by ri.type "
	log.Debugf("select sum(ri.price),ri.type from retail_item ri join retail r on ri.retail_id=r.id and to_char(to_timestamp(extract(epoch from r.create_at )::integer), 'YYYY-MM-DD')='%s' and r.school_id='%s'  group by ri.type ", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
	rows, err := DB.Query(query, goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)

	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	defer rows.Close()
	var oldbookSales, newbookSales, bookType int64
	var price sql.NullInt64
	for rows.Next() {
		err = rows.Scan(&price, &bookType)
		if err != nil {
			log.Error(err)
			return err
		}
		if bookType == 0 {
			if price.Valid {
				newbookSales += price.Int64
			} else {
				newbookSales += 0
			}

		} else {
			if price.Valid {
				oldbookSales += price.Int64
			} else {
				oldbookSales += 0
			}
		}

	}
	goodsSalesStatisticModel.OfflineNewBookSalesFee = newbookSales
	goodsSalesStatisticModel.OfflineOldBookSalesFee = oldbookSales

	return nil
}

//检查学校 在 datetime 是否有数据
func HasThisDayGoodsSalesData(school_id, datetime string) (bool, error) {
	query := "select id from statistic_goods_sales where to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2"
	log.Debugf("select id from statistic_goods_sales where to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s'", datetime, school_id)
	var id string
	err := DB.QueryRow(query, datetime, school_id).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		log.Error(err)
		return false, err
	}
	if id != "" {
		return true, nil
	}
	return false, nil
}

//增肌商品销售数据
func AddGoodsSalesStatistic(model *pb.GoodsSalesStatisticModel, time time.Time) error {
	isExist, err := HasThisDayGoodsSalesData(model.SchoolId, model.StatisticAt)
	if err != nil {
		log.Error(err)
		return err
	}
	if isExist {

		return nil
	}
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()
	query := "insert into statistic_goods_sales (store_id,school_id,alipay_order_num,alipay_order_fee,wechat_order_num,wechat_order_fee,online_new_book_sales_fee,online_old_book_sales_fee,send_order_num,after_sale_num,after_sale_handled_num,after_sale_handled_fee,offline_new_book_sales_fee,offline_old_book_sales_fee,offline_order_num,statistic_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)"
	log.Debugf("insert into statistic_goods_sales (store_id,school_id,alipay_order_num,alipay_order_fee,wechat_order_num,wechat_order_fee,online_new_book_sales_fee,online_old_book_sales_fee,send_order_num,after_sale_num,after_sale_handled_num,after_sale_handled_fee,offline_new_book_sales_fee,offline_old_book_sales_fee,offline_order_num,statistic_at) values('%s','%s',%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,'%+v')", model.StoreId, model.SchoolId, model.AlipayOrderNum, model.AlipayOrderFee, model.WechatOrderNum, model.WechatOrderFee, model.OnlineNewBookSalesFee, model.OnlineOldBookSalesFee, model.SendOrderNum, model.AfterSaleNum, model.AfterSaleHandledNum, model.AfterSaleHandledFee, model.OfflineNewBookSalesFee, model.OfflineOldBookSalesFee, model.OfflineOrderNum, time)
	_, err = tx.Exec(query, model.StoreId, model.SchoolId, model.AlipayOrderNum, model.AlipayOrderFee, model.WechatOrderNum, model.WechatOrderFee, model.OnlineNewBookSalesFee, model.OnlineOldBookSalesFee, model.SendOrderNum, model.AfterSaleNum, model.AfterSaleHandledNum, model.AfterSaleHandledFee, model.OfflineNewBookSalesFee, model.OfflineOldBookSalesFee, model.OfflineOrderNum, time)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = HasThisDayGoodsSalesData(model.SchoolId, model.StatisticAt)

	if err != nil {
		log.Error(err)
		return nil
	}
	tx.Commit()
	return nil
}
