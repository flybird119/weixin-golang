package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	bookDB "github.com/goushuyun/weixin-golang/books/db"
	. "github.com/goushuyun/weixin-golang/db"
	goodsDB "github.com/goushuyun/weixin-golang/goods/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	schoolDB "github.com/goushuyun/weixin-golang/school/db"
	sellerDB "github.com/goushuyun/weixin-golang/seller/db"
	storeDB "github.com/goushuyun/weixin-golang/store/db"
	"github.com/wothing/log"
)

const (
	BasePoundage = 50
)

//提交数据
func OrderSubmit(tx *sql.Tx, carts []*pb.Cart, orderModel *pb.OrderSubmitModel) (order *pb.Order, noStock string, err error) {

	//获取学校的运费
	school, err := schoolDB.GetSchoolById(orderModel.SchoolId)
	if err != nil {
		misc.LogErr(err)
		return nil, "", err
	}
	nowTime := time.Now()
	order = &pb.Order{}
	//首选创建goods，然后创建订单项
	query := "insert into orders (total_fee,freight,user_id,mobile,name,address,remark,store_id,school_id,order_at,goods_fee,withdrawal_fee ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,0,0) returning id"
	log.Debugf(query+"args : %#v", school.ExpressFee, school.ExpressFee, orderModel.UserId, orderModel.Mobile, orderModel.Name, orderModel.Address, orderModel.Remark, orderModel.StoreId, orderModel.SchoolId, nowTime)
	err = tx.QueryRow(query, school.ExpressFee, school.ExpressFee, orderModel.UserId, orderModel.Mobile, orderModel.Name, orderModel.Address, orderModel.Remark, orderModel.StoreId, orderModel.SchoolId, nowTime).Scan(&order.Id)
	if err != nil {
		misc.LogErr(err)
		return nil, "", err
	}
	//遍历carts
	for i := 0; i < len(carts); i++ {
		noStock, err = AddOrderItem(tx, order, carts[i], nowTime)
		if err != nil {
			misc.LogErr(err)
			return nil, "", err
		}
		if noStock != "" {
			return nil, "noStock", nil
		}
	}
	query = "select order_status,total_fee,freight,goods_fee,user_id,mobile,name,address,remark,store_id,school_id,groupon_id from orders where id=$1"
	log.Debugf("select order_status,total_fee,freight,goods_fee,user_id,mobile,name,address,remark,store_id,school_id,groupon_id from orders where id='%s'", order.Id)
	err = tx.QueryRow(query, order.Id).Scan(&order.OrderStatus, &order.TotalFee, &order.Freight, &order.GoodsFee, &order.UserId, &order.Mobile, &order.Name, &order.Address, &order.Remark, &order.StoreId, &order.SchoolId, &order.GrouponId)
	if err != nil {
		misc.LogErr(err)
		return nil, "", err
	}
	return
}

func AddOrderItem(tx *sql.Tx, order *pb.Order, cart *pb.Cart, nowTime time.Time) (noStack string, err error) {
	//减少库存量

	var (
		query      string
		price      int
		amount     int
		is_selling bool
	)
	noStock := "noStock"
	if cart.Type == 0 {
		query = "update goods set new_book_amount=new_book_amount-$1,new_book_sale_amount=new_book_sale_amount+$2 where id=$3 returning new_book_amount,new_book_price,has_new_book"
		log.Debugf("update goods set new_book_amount=new_book_amount-%d,new_book_sale_amount=new_book_sale_amount+%d where id='%s'returning new_book_amount,new_book_price,has_new_book", cart.Amount, cart.Amount, cart.GoodsId)
		err = tx.QueryRow(query, cart.Amount, cart.Amount, cart.GoodsId).Scan(&amount, &price, &is_selling)
		if err != nil {
			misc.LogErr(err)
			return "", err
		}
		if !is_selling || amount < 0 {

			return noStock, nil
		}
	} else {
		query = "update goods set old_book_amount=old_book_amount-$1,old_book_sale_amount=old_book_sale_amount+$2 where id=$3 returning old_book_amount,old_book_price,has_old_book"
		log.Debugf("update goods set old_book_amount=old_book_amount-%d ,old_book_sale_amount=old_book_sale_amount+%s where id='%s'returning old_book_amount,old_book_price,has_old_book", cart.Amount, cart.GoodsId)
		err = tx.QueryRow(query, cart.Amount, cart.Amount, cart.GoodsId).Scan(&amount, &price, &is_selling)
		if err != nil {
			misc.LogErr(err)
			return "", err
		}
		if !is_selling || amount < 0 {

			return noStock, nil
		}
	}
	//然后创建订单项
	query = "insert into orders_item (goods_id,orders_id,type,amount,price,create_at) values($1,$2,$3,$4,$5,$6)"
	log.Debugf("insert into orders_item (goods_id,orders_id,type,amount,price,create_at) values('%s','%s',%d,%d,%d,%v)", cart.GoodsId, order.Id, cart.Type, cart.Amount, price, nowTime)
	_, err = tx.Exec(query, cart.GoodsId, order.Id, cart.Type, cart.Amount, price, nowTime)
	if err != nil {
		misc.LogErr(err)
		return "", err
	}
	//更改订单
	totalFee := int(cart.Amount) * price
	query = "update orders set total_fee=total_fee+$1,goods_fee=goods_fee+$2 where id=$3"
	log.Debugf("update orders set total_fee=total_fee+%d,goods_fee=goods_fee+%d where id='%s'", totalFee, totalFee, order.Id)
	_, err = tx.Exec(query, totalFee, totalFee, order.Id)
	if err != nil {
		misc.LogErr(err)
		return "", err
	}
	return "", nil
}

//订单支付成功
func PaySuccess(order *pb.Order) (isChange bool, err error) {
	isChange = false
	var poundage int64
	poundage = BasePoundage
	query := "select order_status,total_fee, freight,goods_fee,store_id,school_id from orders where id=$1"
	log.Debugf("select order_status,total_fee, freight,goods_fee,store_id,school_id from orders where id=%s", order.Id)
	err = DB.QueryRow(query, order.Id).Scan(&order.OrderStatus, &order.TotalFee, &order.Freight, &order.GoodsFee, &order.StoreId, &order.SchoolId)
	if err != nil {
		misc.LogErr(err)
		return
	}
	if order.OrderStatus != 0 {
		//已经付过款了
		isChange = true
		return
	}
	info := &pb.StoreExtraInfo{StoreId: order.StoreId}
	err = storeDB.GetStoreExtraInfo(info)
	if err != nil {
		log.Error(err)
		return
	}
	if info.Id != "" {
		poundage = info.Poundage
	}
	var serviceDiscount = float64(poundage) / 1000
	//修改订单状态  // 计算待体现金额
	//计算待体现金额
	withdrawalFeeStr := fmt.Sprintf("%0.0f", float64(order.TotalFee)*(1.00-serviceDiscount))
	withdrawalFee, err := strconv.ParseInt(withdrawalFeeStr, 10, 64)
	if err != nil {
		misc.LogErr(err)
		return
	}
	order.WithdrawalFee = withdrawalFee

	//修改的字段有 order_status,withdrawal_fee,trade_no,pay_channel,pay_at,update_at
	tx, err := DB.Begin()
	if err != nil {
		misc.LogErr(err)
		return
	}
	defer tx.Rollback()
	query = "update orders set order_status=$1,withdrawal_fee=$2,trade_no=$3,pay_channel=$4,pay_at=now(),update_at=now() where id=$5"
	log.Debugf("update orders set order_status=%d,withdrawal_fee=%d,trade_no='%s',pay_channel='%s',pay_at=now(),update_at=now() where id='%s'", 1, order.WithdrawalFee, order.TradeNo, order.PayChannel, order.Id)
	_, err = tx.Exec(query, 1, order.WithdrawalFee, order.TradeNo, order.PayChannel, order.Id)
	if err != nil {
		misc.LogErr(err)
		return
	}
	tx.Commit()
	return
}

//搜索订单列表
func FindOrders(order *pb.Order) (details []*pb.OrderDetail, err error, totalcount int64) {
	//需要的项
	var pay_at, deliver_at, print_at, complete_at, after_sale_apply_at, after_sale_end_at, distribute_at, confirm_at, close_at sql.NullString

	selectParam := "o.id,o.order_status,o.total_fee,o.freight,o.goods_fee,o.withdrawal_fee,o.user_id,o.mobile,o.name,o.address,o.remark,o.store_id,o.school_id,o.trade_no,o.pay_channel,extract(epoch from o.order_at)::bigint,extract(epoch from o.pay_at)::bigint,extract(epoch from o.deliver_at)::bigint,extract(epoch from o.print_at)::bigint,extract(epoch from o.complete_at)::bigint,o.print_staff_id,o.deliver_staff_id,o.after_sale_staff_id,extract(epoch from o.after_sale_apply_at)::bigint,extract(epoch from o.after_sale_end_at)::bigint,o.after_sale_status,o.after_sale_trad_no,o.refund_fee,o.groupon_id,extract(epoch from o.update_at)::bigint, extract(epoch from o.distribute_at)::bigint, o.distribute_staff_id, extract(epoch from o.confirm_at)::bigint,extract(epoch from o.close_at)::bigint,apply_refund_fee,seller_remark,seller_remark_type"

	query := fmt.Sprintf("select %s from orders o where 1=1 ", selectParam)
	selectCountQuery := "select count(*) from orders o where 1=1"
	var args []interface{}
	var condition string

	//检索条件
	//1.0 根据状态
	//1.1 根据状态 --> 在根据状态，还要考虑 搜索来源，区别app端搜索和seller端搜索
	if order.OrderStatus != -1 {

		if order.OrderStatus == 80 {
			//1.1.1 如果是app端并且搜索全部的订单，那么显示该用户所有状态的订单
			if order.SearchType != 0 {
				//1.1.3 如果是seller端并且搜索全部的订单，那么显示不为 代付款 ：0 和 已关闭订单：8 状态的订单
				condition += "and (o.order_status <> 0 and o.order_status <> 8) "
			}
			//查看待处理订单
		} else if order.OrderStatus == 79 {
			//所有售后
			condition += fmt.Sprintf(" and o.after_sale_status >= 1")
		} else if order.OrderStatus == 78 {
			//已接受售后
			condition += fmt.Sprintf(" and o.after_sale_status > 1")
		} else if order.OrderStatus == 77 {
			//待处理售后
			condition += fmt.Sprintf(" and o.after_sale_status=1")
		} else {
			args = append(args, order.OrderStatus)
			condition += fmt.Sprintf(" and o.order_status=$%d", len(args))
		}
	}
	//2.0 根据用户
	if order.UserId != "" {
		args = append(args, order.UserId)
		condition += fmt.Sprintf(" and o.user_id=$%d", len(args))
	}
	//3.0 根据云店铺
	if order.StoreId != "" {
		args = append(args, order.StoreId)
		condition += fmt.Sprintf(" and o.store_id=$%d", len(args))
	}
	//4.0 根据学校
	if order.SchoolId != "" {
		args = append(args, order.SchoolId)
		condition += fmt.Sprintf(" and o.school_id=$%d", len(args))
	}
	//5.0 根据付款时间 开始 - 结束
	if order.StartAt != 0 && order.EndAt != 0 {
		args = append(args, order.StartAt)
		condition += fmt.Sprintf(" and extract(epoch from o.pay_at)::bigint between $%d and $%d", len(args), len(args)+1)
		args = append(args, order.EndAt)
	}
	//6.0 订单编号
	if order.Id != "" {
		args = append(args, order.Id)
		condition += fmt.Sprintf(" and o.id=$%d", len(args))
	}
	//7.0 收货人手机号
	if order.Mobile != "" {
		args = append(args, order.Mobile)
		condition += fmt.Sprintf(" and o.mobile=$%d", len(args))
	}
	//8.0 姓名
	if order.Name != "" {
		args = append(args, order.Name)
		condition += fmt.Sprintf(" and o.name=$%d", len(args))
	}
	//9.0 isbn
	if order.Isbn != "" {
		args = append(args, order.Isbn)
		condition += fmt.Sprintf(" and (exists (select * from orders_item oi join  goods g on oi.goods_id=g.id where oi.orders_id=o.id and g.isbn=$%d))", len(args))
	}

	//10.0 班级购id
	if order.GrouponId != "" {
		condition += fmt.Sprintf(" and  o.groupon_id='%s'", order.GrouponId)
	}
	//11.0 商家备注类型
	if order.SellerRemarkType != 0 {
		if order.SellerRemarkType == 79 {
			condition += fmt.Sprintf(" and o.seller_remark_type=0")
		} else if order.SellerRemarkType == 80 {
			condition += fmt.Sprintf(" and o.seller_remark_type<>0")
		} else {
			args = append(args, order.SellerRemarkType)
			condition += fmt.Sprintf(" and o.seller_remark_type=$%d", len(args))
		}

	}
	selectCountQuery += condition

	condition += " order by o.update_at desc"
	if order.Page <= 0 {
		order.Page = 1
	}
	if order.Size <= 0 {
		order.Size = 10
	}
	condition += fmt.Sprintf(" OFFSET %d LIMIT %d ", (order.Page-1)*order.Size, order.Size)
	query += condition
	err = DB.QueryRow(selectCountQuery, args...).Scan(&totalcount)
	if err != nil {
		log.Warn(err)
		misc.LogErr(err)
		return details, nil, totalcount
	}
	if totalcount <= 0 {

		return details, nil, 0
	}
	log.Debugf(query+" args :%#v", args)
	rows, err := DB.Query(query, args...)
	if err != nil && err != sql.ErrNoRows {
		misc.LogErr(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		next := &pb.Order{}
		orderDetail := &pb.OrderDetail{}
		err = rows.Scan(&next.Id, &next.OrderStatus, &next.TotalFee, &next.Freight, &next.GoodsFee, &next.WithdrawalFee, &next.UserId, &next.Mobile, &next.Name, &next.Address, &next.Remark, &next.StoreId, &next.SchoolId, &next.TradeNo, &next.PayChannel, &next.OrderAt, &pay_at, &deliver_at, &print_at, &complete_at, &next.PrintStaffId, &next.DeliverStaffId, &next.AfterSaleStaffId, &after_sale_apply_at, &after_sale_end_at, &next.AfterSaleStatus, &next.AfterSaleTradeNo, &next.RefundFee, &next.GrouponId,
			&next.UpdateAt, &distribute_at, &next.DistributeStaffId, &confirm_at, &close_at, &next.ApplyRefundFee, &next.SellerRemark, &next.SellerRemarkType)
		if err != nil {
			misc.LogErr(err)
			return nil, err, totalcount
		}
		orderitems, err := GetOrderItems(next)
		if err != nil {
			misc.LogErr(err)
			return nil, err, totalcount
		}
		orderDetail.Order = next
		orderDetail.Items = orderitems
		details = append(details, orderDetail)
		//转换可能为空的值
		if pay_at.Valid {
			next.PayAt, _ = strconv.ParseInt(pay_at.String, 10, 64)
		}
		if deliver_at.Valid {
			next.DeliverAt, _ = strconv.ParseInt(deliver_at.String, 10, 64)
		}
		if print_at.Valid {
			next.PrintAt, _ = strconv.ParseInt(print_at.String, 10, 64)
		}
		if complete_at.Valid {
			next.CompleteAt, _ = strconv.ParseInt(complete_at.String, 10, 64)
		}
		if after_sale_apply_at.Valid {
			next.AfterSaleApplyAt, _ = strconv.ParseInt(after_sale_apply_at.String, 10, 64)
		}
		if after_sale_end_at.Valid {
			next.AfterSaleEndAt, _ = strconv.ParseInt(after_sale_end_at.String, 10, 64)
		}
		if distribute_at.Valid {
			next.DistributeAt, _ = strconv.ParseInt(distribute_at.String, 10, 64)
		}
		if confirm_at.Valid {
			next.ConfirmAt, _ = strconv.ParseInt(confirm_at.String, 10, 64)
		}
		if close_at.Valid {
			next.CloseAt, _ = strconv.ParseInt(close_at.String, 10, 64)
		}
	}

	return
}

//根据订单获取订单项集合
func GetOrderItems(order *pb.Order) (orderitems []*pb.OrderItem, err error) {
	query := "select oi.id,g.id,oi.type,oi.amount,oi.price,b.title,b.isbn,b.image,b.price from orders_item oi join goods g on oi.goods_id=g.id join books b on g.book_id=b.id where orders_id='%s'"
	query = fmt.Sprintf(query, order.Id)
	log.Debugf(query)
	rows, err := DB.Query(query)
	//如果出现无结果异常
	if err == sql.ErrNoRows {
		return orderitems, nil
	}
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	//遍历搜索结果
	for rows.Next() {
		orderItem := &pb.OrderItem{}
		orderitems = append(orderitems, orderItem)
		err = rows.Scan(&orderItem.Id, &orderItem.GoodsId, &orderItem.Type, &orderItem.Amount, &orderItem.Price, &orderItem.BookTitle, &orderItem.BookIsbn, &orderItem.BookImage, &orderItem.OriginPrice)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}

		locations, err := goodsDB.GetGoodsLocationDetailByIdAndType(orderItem.GoodsId, orderItem.Type)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		var printLocation string
		for i := 0; i < len(locations); i++ {
			if i == 0 {
				printLocation = locations[i].StorehouseName + "-" + locations[i].ShelfName + "-" + locations[i].FloorName
			} else {
				printLocation += locations[i].StorehouseName + "-" + locations[i].ShelfName + "-" + locations[i].FloorName

			}
		}
		orderItem.PrintLocation = printLocation
	}
	return
}

//更改时间
func UpdateOrder(order *pb.Order) error {
	query := "update orders set id=id"

	var args []interface{}
	var condition string
	if order.UpdateAt != 0 {
		condition += ",update_at=now()"
	}
	if order.OrderStatus != 0 {
		args = append(args, order.OrderStatus)
		condition += fmt.Sprintf(",order_status=order_status|$%d", len(args))
	}
	if order.DeliverAt != 0 {
		args = append(args, time.Now())
		condition += fmt.Sprintf(",deliver_at=$%d", len(args))
	}
	if order.PrintAt != 0 {
		args = append(args, time.Now())
		condition += fmt.Sprintf(",print_at=$%d", len(args))
	}
	if order.CompleteAt != 0 {
		args = append(args, time.Now())
		condition += fmt.Sprintf(",complete_at=$%d", len(args))
	}
	if order.ConfirmAt != 0 {
		args = append(args, time.Now())
		condition += fmt.Sprintf(",confirm_at=$%d", len(args))
	}
	if order.DistributeAt != 0 {
		args = append(args, time.Now())
		condition += fmt.Sprintf(",distribute_at=$%d", len(args))
	}
	if order.PrintStaffId != "" {
		args = append(args, order.PrintStaffId)
		condition += fmt.Sprintf(",print_staff_id=$%d", len(args))
	}
	if order.DeliverStaffId != "" {
		args = append(args, order.DeliverStaffId)
		condition += fmt.Sprintf(",deliver_staff_id=$%d", len(args))
	}
	if order.DistributeStaffId != "" {
		args = append(args, order.DistributeStaffId)
		condition += fmt.Sprintf(",distribute_staff_id=$%d", len(args))
	}
	if order.SellerRemark != "" {
		args = append(args, order.SellerRemark)
		condition += fmt.Sprintf(",seller_remark=$%d", len(args))
	}
	if order.SellerRemarkType != 0 {
		args = append(args, order.SellerRemarkType)
		condition += fmt.Sprintf(",seller_remark_type=$%d", len(args))
	}

	//order_id
	args = append(args, order.Id)
	condition += fmt.Sprintf(" where id=$%d", len(args))
	query += condition
	log.Debugf(query+" args:", args)
	//建立事务
	tx, err := DB.Begin()
	if err != nil {
		misc.LogErr(err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, args...)
	if err != nil {
		misc.LogErr(err)
		return err
	}

	tx.Commit()
	return nil
}

//获取订单的信息
func GetOrderBaseInfo(order *pb.Order) error {
	next := order
	//需要的项
	var pay_at, deliver_at, print_at, complete_at, after_sale_apply_at, after_sale_end_at, distribute_at, confirm_at, close_at sql.NullString

	selectParam := "o.id,o.order_status,o.total_fee,o.freight,o.goods_fee,o.withdrawal_fee,o.user_id,o.mobile,o.name,o.address,o.remark,o.store_id,o.school_id,o.trade_no,o.pay_channel,extract(epoch from o.order_at)::bigint,extract(epoch from o.pay_at)::bigint,extract(epoch from o.deliver_at)::bigint,extract(epoch from o.print_at)::bigint,extract(epoch from o.complete_at)::bigint,o.print_staff_id,o.deliver_staff_id,o.after_sale_staff_id,extract(epoch from o.after_sale_apply_at)::bigint,extract(epoch from o.after_sale_end_at)::bigint,o.after_sale_status,o.after_sale_trad_no,o.refund_fee,o.groupon_id,extract(epoch from o.update_at)::bigint,extract(epoch from o.distribute_at)::bigint,distribute_staff_id,extract(epoch from o.confirm_at)::bigint, extract(epoch from o.close_at)::bigint,o.apply_refund_fee,seller_remark,seller_remark_type"

	query := fmt.Sprintf("select %s from orders o where id=$1 ", selectParam)

	err := DB.QueryRow(query, order.Id).Scan(&next.Id, &next.OrderStatus, &next.TotalFee, &next.Freight, &next.GoodsFee, &next.WithdrawalFee, &next.UserId, &next.Mobile, &next.Name, &next.Address, &next.Remark, &next.StoreId, &next.SchoolId, &next.TradeNo, &next.PayChannel, &next.OrderAt, &pay_at, &deliver_at, &print_at, &complete_at, &next.PrintStaffId, &next.DeliverStaffId, &next.AfterSaleStaffId, &after_sale_apply_at, &after_sale_end_at, &next.AfterSaleStatus, &next.AfterSaleTradeNo, &next.RefundFee, &next.GrouponId,
		&next.UpdateAt, &distribute_at, &next.DistributeStaffId, &confirm_at, &close_at, &next.ApplyRefundFee, &next.SellerRemark, &next.SellerRemarkType)
	if err != nil && err != sql.ErrNoRows {
		misc.LogErr(err)
		return err
	}
	//转换可能为空的值
	if pay_at.Valid {
		next.PayAt, _ = strconv.ParseInt(pay_at.String, 10, 64)
	}
	if deliver_at.Valid {
		next.DeliverAt, _ = strconv.ParseInt(deliver_at.String, 10, 64)
	}
	if print_at.Valid {
		next.PrintAt, _ = strconv.ParseInt(print_at.String, 10, 64)
	}
	if complete_at.Valid {
		next.CompleteAt, _ = strconv.ParseInt(complete_at.String, 10, 64)
	}
	if after_sale_apply_at.Valid {
		next.AfterSaleApplyAt, _ = strconv.ParseInt(after_sale_apply_at.String, 10, 64)
	}
	if after_sale_end_at.Valid {
		next.AfterSaleEndAt, _ = strconv.ParseInt(after_sale_end_at.String, 10, 64)
	}
	if distribute_at.Valid {
		next.DistributeAt, _ = strconv.ParseInt(distribute_at.String, 10, 64)
	}
	if confirm_at.Valid {
		next.ConfirmAt, _ = strconv.ParseInt(confirm_at.String, 10, 64)
	}
	if close_at.Valid {
		next.CloseAt, _ = strconv.ParseInt(close_at.String, 10, 64)
	}

	return nil
}

//获取订单的信息
func GetOrderBaseInfoByTradeNo(order *pb.Order) error {
	next := order
	//需要的项
	var pay_at, deliver_at, print_at, complete_at, after_sale_apply_at, after_sale_end_at, distribute_at, confirm_at, close_at sql.NullString

	selectParam := "o.id,o.order_status,o.total_fee,o.freight,o.goods_fee,o.withdrawal_fee,o.user_id,o.mobile,o.name,o.address,o.remark,o.store_id,o.school_id,o.trade_no,o.pay_channel,extract(epoch from o.order_at)::bigint,extract(epoch from o.pay_at)::bigint,extract(epoch from o.deliver_at)::bigint,extract(epoch from o.print_at)::bigint,extract(epoch from o.complete_at)::bigint,o.print_staff_id,o.deliver_staff_id,o.after_sale_staff_id,extract(epoch from o.after_sale_apply_at)::bigint,extract(epoch from o.after_sale_end_at)::bigint,o.after_sale_status,o.after_sale_trad_no,o.refund_fee,o.groupon_id,extract(epoch from o.update_at)::bigint,extract(epoch from o.distribute_at)::bigint,distribute_staff_id,extract(epoch from o.confirm_at)::bigint, extract(epoch from o.close_at)::bigint,o.apply_refund_fee"

	query := fmt.Sprintf("select %s from orders o where trade_no=$1 ", selectParam)

	err := DB.QueryRow(query, order.TradeNo).Scan(&next.Id, &next.OrderStatus, &next.TotalFee, &next.Freight, &next.GoodsFee, &next.WithdrawalFee, &next.UserId, &next.Mobile, &next.Name, &next.Address, &next.Remark, &next.StoreId, &next.SchoolId, &next.TradeNo, &next.PayChannel, &next.OrderAt, &pay_at, &deliver_at, &print_at, &complete_at, &next.PrintStaffId, &next.DeliverStaffId, &next.AfterSaleStaffId, &after_sale_apply_at, &after_sale_end_at, &next.AfterSaleStatus, &next.AfterSaleTradeNo, &next.RefundFee, &next.GrouponId,
		&next.UpdateAt, &distribute_at, &next.DistributeStaffId, &confirm_at, &close_at, &next.ApplyRefundFee)
	if err != nil && err != sql.ErrNoRows {
		misc.LogErr(err)
		return err
	}
	//转换可能为空的值
	if pay_at.Valid {
		next.PayAt, _ = strconv.ParseInt(pay_at.String, 10, 64)
	}
	if deliver_at.Valid {
		next.DeliverAt, _ = strconv.ParseInt(deliver_at.String, 10, 64)
	}
	if print_at.Valid {
		next.PrintAt, _ = strconv.ParseInt(print_at.String, 10, 64)
	}
	if complete_at.Valid {
		next.CompleteAt, _ = strconv.ParseInt(complete_at.String, 10, 64)
	}
	if after_sale_apply_at.Valid {
		next.AfterSaleApplyAt, _ = strconv.ParseInt(after_sale_apply_at.String, 10, 64)
	}
	if after_sale_end_at.Valid {
		next.AfterSaleEndAt, _ = strconv.ParseInt(after_sale_end_at.String, 10, 64)
	}
	if distribute_at.Valid {
		next.DistributeAt, _ = strconv.ParseInt(distribute_at.String, 10, 64)
	}
	if confirm_at.Valid {
		next.ConfirmAt, _ = strconv.ParseInt(confirm_at.String, 10, 64)
	}
	if close_at.Valid {
		next.CloseAt, _ = strconv.ParseInt(close_at.String, 10, 64)
	}

	return nil
}

//获取订单员工信息
func GetOrderStaffWork(order *pb.Order) (staffs []*pb.OrderStaff, err error) {
	//print
	if order.PrintStaffId != "" {
		staff, err := getStaffInfo(staffs, order.PrintStaffId, "print")
		if err != nil {
			log.Debug(err)
			misc.LogErr(err)
			return nil, err
		}
		staffs = append(staffs, staff)
	}
	//deliver
	if order.DeliverStaffId != "" {
		staff, err := getStaffInfo(staffs, order.DeliverStaffId, "deliver")
		if err != nil {
			log.Debug(err)
			misc.LogErr(err)
			return nil, err
		}
		staffs = append(staffs, staff)
	}
	//distribute
	if order.DistributeStaffId != "" {
		staff, err := getStaffInfo(staffs, order.DistributeStaffId, "distribute")
		if err != nil {
			log.Debug(err)
			misc.LogErr(err)
			return nil, err
		}
		staffs = append(staffs, staff)
	}
	//after_sale
	if order.AfterSaleStaffId != "" {
		staff, err := getStaffInfo(staffs, order.AfterSaleStaffId, "after_sale")
		if err != nil {
			log.Debug(err)
			misc.LogErr(err)
			return nil, err
		}
		staffs = append(staffs, staff)
	}
	return staffs, nil
}

//获取售后详情
func GetAfterSaleDetail(order *pb.Order) (*pb.AfterSaleModel, error) {
	log.Debug(order)
	if order.AfterSaleStatus != 0 {

		query := "select o.refund_fee ,o.after_sale_reason,o.after_sale_images,o.apply_refund_fee,o.after_sale_trad_no from orders o where o.id=$1"
		log.Debugf("select o.refund_fee ,o.after_sale_reason,o.after_sale_images,o.apply_refund_fee,o.after_sale_trad_no from orders where o.id='%s'", order.Id)
		var images []*pb.AfterSaleImage
		var imageStr string
		var after_sale_trad_no sql.NullString
		afterSaleModdel := &pb.AfterSaleModel{}
		err := DB.QueryRow(query, order.Id).Scan(&afterSaleModdel.RefundFee, &afterSaleModdel.Reason, &imageStr, &afterSaleModdel.ApplyRefundFee, &after_sale_trad_no)
		if err != nil {
			log.Debug(err)
			misc.LogErr(err)
			return nil, err
		}
		if after_sale_trad_no.Valid {
			afterSaleModdel.RefundTradeNo = after_sale_trad_no.String
		}
		//转化staff
		if err := json.Unmarshal([]byte(imageStr), &images); err != nil {
			log.Debug(err)
			misc.LogErr(err)
			return nil, err
		}
		afterSaleModdel.Images = images

		return afterSaleModdel, nil
	}
	return nil, nil
}

//私有方法 获取用户信息
func getStaffInfo(staffs []*pb.OrderStaff, findStaffId string, workName string) (*pb.OrderStaff, error) {

	for index := 0; index < len(staffs); index++ {
		staff := staffs[index]
		if staff.StaffId == findStaffId {
			findStaff := &pb.OrderStaff{StaffId: findStaffId, StaffName: staff.StaffName, StaffWork: workName}

			return findStaff, nil
		}
	}
	sellerInfo, err := sellerDB.GetSellerById(findStaffId)
	if err != nil {
		log.Debug(err)
		misc.LogErr(err)
		return nil, err
	}
	findStaff := &pb.OrderStaff{StaffId: sellerInfo.Id, StaffName: sellerInfo.Username, StaffWork: workName}

	return findStaff, nil
}

//填写售后信息
func FillInAfterSaleDetail(aftersaleModel *pb.AfterSaleModel) error {
	query := "update orders set after_sale_apply_at=now(),update_at=now()"

	var args []interface{}
	var condition string
	//退款原因
	if aftersaleModel.Reason != "" {
		args = append(args, aftersaleModel.Reason)
		condition += fmt.Sprintf(",after_sale_reason=$%d", len(args))
	}

	imagesStr, err := json.Marshal(aftersaleModel.Images)
	if err == nil {
		args = append(args, imagesStr)
		condition += fmt.Sprintf(",after_sale_images=$%d", len(args))
	}
	//退款费用
	args = append(args, aftersaleModel.RefundFee)
	condition += fmt.Sprintf(",refund_fee=$%d,apply_refund_fee=$%d", len(args), len(args)+1)
	args = append(args, aftersaleModel.RefundFee)
	//订单退款状态
	args = append(args, 1)
	condition += fmt.Sprintf(",after_sale_status=$%d", len(args))
	//订单状态
	args = append(args, 16)
	condition += fmt.Sprintf(",order_status=order_status|$%d", len(args))

	args = append(args, aftersaleModel.OrderId)
	condition += fmt.Sprintf(" where id=$%d", len(args))

	args = append(args, aftersaleModel.UserId)
	condition += fmt.Sprintf(" and user_id=$%d ", len(args))

	query += condition
	log.Debugf(query+" args:%+v", args)
	_, err = DB.Exec(query, args...)
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	return nil
}

//关闭订单处理
func CloseOrder(order *pb.Order) error {
	//首先更改订单的状态
	query := "update orders set order_status=order_status|$1,close_at=now() where id=$2 and user_id=$3 returning order_status"
	log.Debugf("update orders set order_status=order_status|%d,close_at=now() where id='%s' and user_id='%s' returning order_status", 8, order.Id, order.UserId)
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, 8, order.Id, order.UserId).Scan(&order.OrderStatus)
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	if order.OrderStatus != 8 {
		return errors.New("order's state err!")
	}
	//获取订单的订单项
	orderitems, err := GetOrderItems(order)
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	for i := 0; i < len(orderitems); i++ {
		//恢复商品的库存
		orderitem := orderitems[i]
		err := goodsDB.RecoverGoodsAmountFromClosedOrder(tx, orderitem)
		if err != nil {
			log.Error(err)
			misc.LogErr(err)
			return err
		}
	}
	tx.Commit()
	return nil
}

//给订单打备注
func RemarkOrder(order *pb.Order) error {
	query := "update orders set seller_remark='%s',seller_remark_type=%d where id='%s'"
	query = fmt.Sprintf(query, order.SellerRemark, order.SellerRemarkType, order.Id)
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//处理售后订单
func HandleAfterSaleOrder(tx *sql.Tx, order *pb.Order) error {
	//修改状态，退款金额
	query := "update orders set update_at=now()"

	var condition string

	//区分退款金额
	if order.RefundFee == 0 {
		condition += fmt.Sprintf(",refund_fee=%d,after_sale_status=4,after_sale_end_at=now()", order.RefundFee)
	} else {
		condition += fmt.Sprintf(",refund_fee=%d,after_sale_status=2", order.RefundFee)
	}

	condition += fmt.Sprintf(",after_sale_staff_id='%s'", order.AfterSaleStaffId)
	condition += fmt.Sprintf(" where id='%s' returning after_sale_status", order.Id)
	query += condition
	log.Debugf(query)
	err := tx.QueryRow(query).Scan(&order.AfterSaleStatus)
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	return nil
}

//用户个人中心必要订单数量统计
func UserCenterNecessaryOrderCount(model *pb.UserCenterOrderCount) error {
	query := "select count(*) from orders where order_status=0 and user_id=$1 and store_id=$2"
	log.Debugf("select count(*) from orders where order_status=0 and user_id='%s' and store_id='%s'", model.UserId, model.StoreId)
	err := DB.QueryRow(query, model.UserId, model.StoreId).Scan(&model.UnpaidOrderNum)
	if err != nil {
		log.Error(err)
		return err
	}
	query = "select count(*) from orders where order_status=1 and user_id=$1 and store_id=$2"
	log.Debugf("select count(*) from orders where order_status=1 and user_id='%s' and store_id='%s'", model.UserId, model.StoreId)
	err = DB.QueryRow(query, model.UserId, model.StoreId).Scan(&model.UndeliveredOrderNum)
	if err != nil {
		log.Error(err)
		return err
	}
	query = "select count(*) from orders where order_status=3 and user_id=$1 and store_id=$2"
	log.Debugf("select count(*) from orders where order_status=3 and user_id='%s' and store_id='%s'", model.UserId, model.StoreId)
	err = DB.QueryRow(query, model.UserId, model.StoreId).Scan(&model.UncompletedOrderNum)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//店铺首页必要订单状态数量
func StoreCenterNecessaryOrderCount(model *pb.StoreHistoryStateOrderNumModel) error {
	//统计待发货订单量
	query := fmt.Sprintf("select count(*) from orders where order_status=1 and store_id='%s'", model.StoreId)
	if model.SchoolId != "" {
		query += fmt.Sprintf(" and school_id='%s'", model.SchoolId)
	}
	log.Debugf(query)
	err := DB.QueryRow(query).Scan(&model.UndeliveredOrderNum)
	if err != nil {
		log.Error(err)
		return err
	}
	//统计待售后订单量
	query = fmt.Sprintf("select count(*) from orders where after_sale_status=1 and store_id='%s'", model.StoreId)
	if model.SchoolId != "" {
		query += fmt.Sprintf(" and school_id='%s'", model.SchoolId)
	}
	log.Debugf(query)
	err = DB.QueryRow(query).Scan(&model.AfterSaleOrderNum)
	if err != nil {
		log.Error(err)
		return err
	}
	todayStatistic := &pb.StoreStatisticDaliyModel{}
	yesterdayStatistic := &pb.StoreStatisticDaliyModel{}
	now := time.Now()
	todayDateStr := now.Format("2006-01-02")
	now = now.Add(-1 * 24 * time.Hour)
	yesterDayDateStr := now.Format("2006-01-02")
	//统计今日处理订单量-以及订单金额
	query = fmt.Sprintf("select count(*),sum(total_fee),to_char(to_timestamp(extract(epoch from pay_at)::bigint), 'YYYY-MM-DD') as d from orders where to_char(to_timestamp(extract(epoch from pay_at)::bigint), 'YYYY-MM-DD') in ('%s','%s') and store_id='%s'", yesterDayDateStr, todayDateStr, model.StoreId)
	if model.SchoolId != "" {
		query += fmt.Sprintf(" and school_id='%s'", model.SchoolId)
	}
	query += " group by d"
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	defer rows.Close()
	var totalSales, orderNum int64
	var dateStr string
	for rows.Next() {
		err = rows.Scan(&orderNum, &totalSales, &dateStr)
		if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			return err
		}
		if dateStr == todayDateStr {
			todayStatistic.TotalSales = totalSales
			todayStatistic.OrderNum = orderNum
		} else {
			yesterdayStatistic.TotalSales = totalSales
			yesterdayStatistic.OrderNum = orderNum
		}
	}

	//统计处理的订单量
	query = fmt.Sprintf("select count(*),to_char(to_timestamp(extract(epoch from pay_at)::bigint), 'YYYY-MM-DD') as d from orders where to_char(to_timestamp(extract(epoch from deliver_at)::bigint), 'YYYY-MM-DD') in ('%s','%s') and store_id='%s'", yesterDayDateStr, todayDateStr, model.StoreId)
	if model.SchoolId != "" {
		query += fmt.Sprintf(" and school_id='%s'", model.SchoolId)
	}
	query += " group by d"
	log.Debug(query)
	rows, err = DB.Query(query)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&orderNum, &dateStr)
		if err != nil && err != sql.ErrNoRows {
			log.Error(err)
			return err
		}
		if dateStr == todayDateStr {
			todayStatistic.HandledOrderNum = orderNum
		} else {
			yesterdayStatistic.HandledOrderNum = orderNum
		}
	}
	model.TodayData = todayStatistic
	model.YesterdayData = yesterdayStatistic
	return nil
}

//售后结果处理
func AfterSaleResultOperation(afterSaleModel *pb.AfterSaleModel) error {
	query := "update orders set after_sale_status=%d,after_sale_trad_no=%s,refund_fee=%d,after_sale_end_at=now() where trade_no='%s' returning id"
	if afterSaleModel.IsSuccess {
		query = fmt.Sprintf(query, 4, afterSaleModel.RefundTradeNo, afterSaleModel.RefundFee, afterSaleModel.TradeNo)
	} else {
		query = fmt.Sprintf(query, 3, afterSaleModel.RefundTradeNo, afterSaleModel.RefundFee, afterSaleModel.TradeNo)
	}
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//导出订单data
func ExportDeliveryOrderData(order *pb.Order) (details []*pb.OrderDetail, err error) {
	//需要的项
	//导出代打印订单
	query := fmt.Sprintf("select id,mobile,name,address,remark from orders o where order_status=1")
	var args []interface{}
	var condition string

	//检索条件
	//检索条件 -- 云店铺
	if order.StoreId != "" {
		args = append(args, order.StoreId)
		condition += fmt.Sprintf(" and o.store_id=$%d", len(args))
	}
	//检索条件 -- 订单ids
	condition += fmt.Sprintf(" and o.id in(%s)  order by o.update_at desc", order.Ids)

	query += condition
	log.Debugf(query+" args :%#v", args)
	rows, err := DB.Query(query, args...)
	if err != nil && err != sql.ErrNoRows {
		misc.LogErr(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		next := &pb.Order{}
		orderDetail := &pb.OrderDetail{}
		//select id,mobile,name,address,remark
		err = rows.Scan(&next.Id, &next.Mobile, &next.Name, &next.Address, &next.Remark)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		orderitems, err := getExportOrderItems(next)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		orderDetail.Order = next
		orderDetail.Items = orderitems
		details = append(details, orderDetail)

	}

	return
}

//导出配货单
func ExportDistributeOrderData(order *pb.Order) (models []*pb.DistributeOrderModel, err error) {
	query := "select goods_id,type,sum(amount) from orders_item where orders_id in(%s) group by goods_id,type order by goods_id desc"
	query = fmt.Sprintf(query, order.Ids)
	log.Debug(query)

	rows, err := DB.Query(query)
	if err != nil && err != sql.ErrNoRows {
		misc.LogErr(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		next := &pb.DistributeOrderModel{}
		models = append(models, next)
		var goodsId string
		err = rows.Scan(&goodsId, &next.Type, &next.Num)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		//获取商品信息
		goods := &pb.Goods{Id: goodsId, StoreId: order.StoreId}
		err = goodsDB.GetGoodsByIdOrIsbn(goods)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		//获取商品所对应书本详情
		book := &pb.Book{Id: goods.BookId}
		err = bookDB.GetBookInfo(book)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		next.Isbn = book.Isbn
		next.Title = book.Title
		next.Publisher = book.Publisher
		//获取商品位置
		locations, err := goodsDB.GetGoodsLocationDetailByIdAndType(goodsId, next.Type)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		var locationStr string
		for i := 0; i < len(locations); i++ {
			if i == 0 {
				locationStr = locations[i].StorehouseName + "--" + locations[i].ShelfName + "--" + locations[i].FloorName
			} else {
				locationStr += locations[i].StorehouseName + "--" + locations[i].ShelfName + "--" + locations[i].FloorName
				locationStr += "，"
			}
		}
		next.Locations = locationStr
	}
	return
}

//导出订单获取订单项详情
func getExportOrderItems(order *pb.Order) (orderitems []*pb.OrderItem, err error) {
	query := "select oi.id,g.id,oi.type,oi.amount,oi.price,b.title,b.isbn,b.image,b.price from orders_item oi join goods g on oi.goods_id=g.id join books b on g.book_id=b.id where orders_id='%s'"
	query = fmt.Sprintf(query, order.Id)
	log.Debugf(query)
	rows, err := DB.Query(query)

	//如果出现无结果异常
	if err == sql.ErrNoRows {
		return orderitems, nil
	}
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	//遍历搜索结果
	for rows.Next() {
		orderItem := &pb.OrderItem{}
		orderitems = append(orderitems, orderItem)
		err = rows.Scan(&orderItem.Id, &orderItem.GoodsId, &orderItem.Type, &orderItem.Amount, &orderItem.Price, &orderItem.BookTitle, &orderItem.BookIsbn, &orderItem.BookImage, &orderItem.OriginPrice)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		locations, err := goodsDB.GetGoodsLocationDetailByIdAndType(orderItem.GoodsId, orderItem.Type)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		orderItem.Locations = locations
	}
	return
}

func RestatisticOrderNum() error {
	query := "select id,school_id,to_char(to_timestamp(extract(epoch from statistic_at)::bigint), 'YYYY-MM-DD') from statistic_goods_sales"
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Error(err)
		return err
	}
	defer rows.Next()
	for rows.Next() {
		statisticModel := &pb.GoodsSalesStatisticModel{}
		var dateStr string
		err := rows.Scan(&statisticModel.Id, &statisticModel.SchoolId, &dateStr)
		if err != nil {
			log.Error(err)
			return err
		}
		updateStatisticOrderNum(statisticModel, dateStr)
	}
	return nil
}

func updateStatisticOrderNum(model *pb.GoodsSalesStatisticModel, dateStr string) error {
	query := "select count(*),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::bigint), 'YYYY-MM-DD')=$1 and school_id=$2   group by pay_channel"
	log.Debugf("select count(*),pay_channel from orders where to_char(to_timestamp(extract(epoch from pay_at )::bigint), 'YYYY-MM-DD')='%s' and school_id='%s'   group by pay_channel", dateStr, model.SchoolId)
	rows, err := DB.Query(query, dateStr, model.SchoolId)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		return err
	}
	defer rows.Close()
	var alipay_order_num, wechat_order_num, count int64
	var pay_channel string
	for rows.Next() {
		err = rows.Scan(&count, &pay_channel)
		if err != nil {
			log.Error(err)
			return err
		}
		if strings.Contains(pay_channel, "alipay") {
			alipay_order_num += count

		} else {
			wechat_order_num += count
		}
	}

	query = "update statistic_goods_sales set alipay_order_num=%d ,wechat_order_num=%d where id='%s'"
	query = fmt.Sprintf(query, alipay_order_num, wechat_order_num, model.Id)
	log.Debug(query)
	_, err = DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
