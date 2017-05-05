package db

import (
	"database/sql"
	"fmt"
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
	query := "select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2   group by pay_channel"
	log.Debugf("select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s'   group by pay_channel", goodsSalesStatisticModel.StatisticAt, goodsSalesStatisticModel.SchoolId)
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

//某天销售额
func GetOneDaySales(model *pb.GoodsSalesStatisticModel) error {
	var rows *sql.Rows
	var err error

	if model.SchoolId != "" {
		query := "select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::integer), 'YYYY-MM-DD')=$1 and school_id=$2 and store_id=$3  group by pay_channel"
		log.Debugf("select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::integer), 'YYYY-MM-DD')='%s' and school_id='%s' and store_id=$3  group by pay_channel", model.StatisticAt, model.SchoolId, model.StoreId)
		rows, err = DB.Query(query, model.StatisticAt, model.SchoolId, model.StoreId)
	} else {
		query := "select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::integer), 'YYYY-MM-DD')=$1 and store_id=$2   group by pay_channel"
		log.Debugf("select sum(total_fee),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::integer), 'YYYY-MM-DD')='%s' and store_id='%s'   group by pay_channel", model.StatisticAt, model.StoreId)
		rows, err = DB.Query(query, model.StatisticAt, model.StoreId)
	}

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
	model.AlipayOrderNum = alipay_order_num
	model.AlipayOrderFee = alipay_order_fee
	model.WechatOrderNum = wechat_order_num
	model.WechatOrderFee = wechat_order_fee

	//统计新书旧书销售额
	if model.SchoolId != "" {
		query := "select sum(ri.price),ri.type from retail_item ri join retail r on ri.retail_id=r.id and to_char(to_timestamp(extract(epoch from r.create_at )::integer), 'YYYY-MM-DD')=$1 and r.school_id=$2 and r.store_id=$3 group by ri.type "
		log.Debugf("select sum(ri.price),ri.type from retail_item ri join retail r on ri.retail_id=r.id and to_char(to_timestamp(extract(epoch from r.create_at )::integer), 'YYYY-MM-DD')='%s' and r.school_id='%s' and r.store_id='%s'  group by ri.type ", model.StatisticAt, model.SchoolId, model.StoreId)
		rows, err = DB.Query(query, model.StatisticAt, model.SchoolId, model.StoreId)
	} else {
		query := "select sum(ri.price),ri.type from retail_item ri join retail r on ri.retail_id=r.id and to_char(to_timestamp(extract(epoch from r.create_at )::integer), 'YYYY-MM-DD')=$1 and r.school_id=$2  group by ri.type "
		log.Debugf("select sum(ri.price),ri.type from retail_item ri join retail r on ri.retail_id=r.id and to_char(to_timestamp(extract(epoch from r.create_at )::integer), 'YYYY-MM-DD')='%s' and r.store_id='%s'  group by ri.type ", model.StatisticAt, model.StoreId)
		rows, err = DB.Query(query, model.StatisticAt, model.StoreId)
	}
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
	model.OfflineNewBookSalesFee = newbookSales
	model.OfflineOldBookSalesFee = oldbookSales
	return nil
}

//历史销售额
func HistoryTotalSales(model *pb.GoodsSalesStatisticModel) (totalModel *pb.StatisticTotalModel, err error) {
	query := "select sum(alipay_order_fee+wechat_order_fee),sum(offline_new_book_sales_fee+offline_old_book_sales_fee) ,sum(online_new_book_sales_fee+offline_new_book_sales_fee),sum(online_old_book_sales_fee+offline_old_book_sales_fee) from statistic_goods_sales where 1=1"
	if model.SchoolId != "" {
		query = fmt.Sprintf(query+" and school_id='%s'", model.SchoolId)
	}
	query = fmt.Sprintf(query+" and store_id='%s'", model.StoreId)
	var online_total_sales, offline_total_sales, newbook_total_sales, oldbook_total_sales sql.NullInt64
	log.Debugf(query)
	err = DB.QueryRow(query).Scan(&online_total_sales, &offline_total_sales, &newbook_total_sales, &oldbook_total_sales)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if online_total_sales.Valid {
		totalModel.OnlineTotalSales = online_total_sales.Int64
	}
	if offline_total_sales.Valid {
		totalModel.OfflineTotalSales = offline_total_sales.Int64
	}
	if newbook_total_sales.Valid {
		totalModel.NewbookTotalSales = newbook_total_sales.Int64
	}
	if oldbook_total_sales.Valid {
		totalModel.OldbookTotalSales = oldbook_total_sales.Int64
	}
	return totalModel, nil
}

func HistoryDaliySales(model *pb.GoodsSalesStatisticModel) (salesModels []*pb.GoodsSalesStatisticModel, err error) {
	//time.Unix(time.Now().Unix(), 0)
	var startAt, endAt string
	now := time.Now()
	if model.StartAt == 0 || model.EndAt == 0 {
		startAt = (now.Add(-1 * 24 * time.Hour)).Format("2006-01-02")
		endAt = (now.Add(-15 * 24 * time.Hour)).Format("2006-01-02")

	} else {
		startAt = time.Unix(model.StartAt, 0).Format("2006-01-02")

		endAt = time.Unix(model.EndAt, 0).Format("2006-01-02")
	}
	log.Debugf("start_at:%s and end_at:%s", startAt, endAt)

	query := "select alipay_order_num,alipay_order_fee,wechat_order_num,wechat_order_fee,online_new_book_sales_fee,online_old_book_sales_fee,send_order_num,after_sale_num,after_sale_handled_num,after_sale_handled_fee,offline_new_book_sales_fee,offline_old_book_sales_fee,offline_order_num,to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM-DD') from statistic_goods_sales where 1=1"
	var condition string
	//拼接字符串
	if model.SchoolId != "" {
		condition += fmt.Sprintf(" and school_id='%s'", model.SchoolId)
	}
	condition += fmt.Sprintf(" and store_id='%s' and to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM-DD') between '%s' and '%s'  order by statistic_at desc", model.StoreId, startAt, endAt)

	query += condition

	log.Debugf(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return salesModels, nil
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		find := &pb.GoodsSalesStatisticModel{}
		salesModels = append(salesModels, find)
		//alipay_order_num,alipay_order_fee,wechat_order_num,wechat_order_fee,online_new_book_sales_fee,online_old_book_sales_fee,send_order_num,after_sale_num,after_sale_handled_num,after_sale_handled_fee,offline_new_book_sales_fee,offline_old_book_sales_fee,offline_order_num,to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM-DD')
		err = rows.Scan(&find.AlipayOrderNum, &find.AlipayOrderFee, &find.WechatOrderNum, &find.WechatOrderFee, &find.OnlineNewBookSalesFee, &find.OnlineOldBookSalesFee, &find.SendOrderNum, &find.AfterSaleNum, &find.AfterSaleHandledNum, &find.AfterSaleHandledFee, &find.OfflineNewBookSalesFee, &find.OfflineOldBookSalesFee, &find.OfflineOrderNum, &find.StatisticAt)
		if err != nil {
			return
		}
	}
	return
}

//月份销售额
func HistoryMonthSales(model *pb.GoodsSalesStatisticModel) (salesModels []*pb.StatisticMonthModel, err error) {
	var startAt, endAt string
	now := time.Now()
	if model.StartAt == 0 || model.EndAt == 0 {
		startAt = (now.AddDate(0, -1, 0)).Format("2006-01")
		endAt = (now.AddDate(0, -7, 0)).Format("2006-01")

	} else {
		startAt = time.Unix(model.StartAt, 0).Format("2006-01")

		endAt = time.Unix(model.EndAt, 0).Format("2006-01")
	}
	log.Debugf("start_at:%s and end_at:%s", startAt, endAt)

	query := "select sum(online_new_book_sales_fee+offline_new_book_sales_fee),sum(online_old_book_sales_fee+offline_old_book_sales_fee),sum(alipay_order_fee+wechat_order_fee),sum(offline_new_book_sales_fee+offline_old_book_sales_fee),to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM') from statistic_goods_sales where 1=1"
	var condition string
	//拼接字符串
	if model.SchoolId != "" {
		condition += fmt.Sprintf(" and school_id='%s'", model.SchoolId)
	}
	condition += fmt.Sprintf(" and store_id='%s' and to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM') between '%s' and '%s' group by to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM')  order by statistic_at desc", model.StoreId, startAt, endAt)

	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return salesModels, nil
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		find := &pb.StatisticMonthModel{}
		salesModels = append(salesModels, find)
		var newbook_sales, oldbook_sales, online_sales, offline_sales sql.NullFloat64
		var month sql.NullString
		//sum(online_new_book_sales_fee+offline_new_book_sales_fee),sum(online_old_book_sales_fee+offline_old_book_sales_fee),sum(alipay_order_fee+wechat_order_fee),sum(offline_new_book_sales_fee+offline_old_book_sales_fee),to_char(to_timestamp(extract(epoch from statistic_at )::integer), 'YYYY-MM')
		err = rows.Scan(&newbook_sales, &oldbook_sales, &online_sales, &offline_sales, &month)
		if err != nil {
			log.Error(err)
			return
		}
		if newbook_sales.Valid {
			find.NewbookSales = int64(newbook_sales.Float64)
		}
		if oldbook_sales.Valid {
			find.OldbookSales = int64(oldbook_sales.Float64)
		}
		if online_sales.Valid {
			find.OnlineSales = int64(oldbook_sales.Float64)
		}
		if offline_sales.Valid {
			find.OfflineSales = int64(offline_sales.Float64)
		}
		if month.Valid {
			find.Month = month.String
		}
	}
	return
}
