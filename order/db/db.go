package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	schoolDB "github.com/goushuyun/weixin-golang/school/db"
	"github.com/wothing/log"
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
		query = "update goods set new_book_amount=new_book_amount-$1 where id=$2 returning new_book_amount,new_book_price,has_new_book"
		log.Debugf("update goods set new_book_amount=new_book_amount-%d where id='%s'returning new_book_amount,new_book_price,has_new_book", cart.Amount, cart.GoodsId)
		err = tx.QueryRow(query, cart.Amount, cart.GoodsId).Scan(&amount, &price, &is_selling)
		if err != nil {
			misc.LogErr(err)
			return "", err
		}
		if !is_selling || amount < 0 {

			return noStock, nil
		}
	} else {
		query = "update goods set old_book_amount=old_book_amount-$1 where id=$2 returning old_book_amount,old_book_price,has_old_book"
		log.Debugf("update goods set old_book_amount=old_book_amount-%d where id='%s'returning old_book_amount,old_book_price,has_old_book", cart.Amount, cart.GoodsId)
		err = tx.QueryRow(query, cart.Amount, cart.GoodsId).Scan(&amount, &price, &is_selling)
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
	var serviceDiscount = 0.02
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
func FindOrders(order *pb.Order) (details []*pb.OrderDetail, err error) {
	//需要的项
	var pay_at, deliver_at, print_at, complete_at, after_sale_apply_at, after_sale_end_at sql.NullString

	selectParam := "o.id,o.order_status,o.total_fee,o.freight,o.goods_fee,o.withdrawal_fee,o.user_id,o.mobile,o.name,o.address,o.remark,o.store_id,o.school_id,o.trade_no,o.pay_channel,extract(epoch from o.order_at)::integer,extract(epoch from o.pay_at)::integer,extract(epoch from o.deliver_at)::integer,extract(epoch from o.print_at)::integer,extract(epoch from o.complete_at)::integer,o.print_staff_id,o.deliver_staff_id,o.after_sale_staff_id,extract(epoch from o.after_sale_apply_at)::integer,extract(epoch from o.after_sale_end_at)::integer,o.after_sale_status,o.after_sale_trad_no,o.refund_fee,o.groupon_id,extract(epoch from o.update_at)::integer"

	query := fmt.Sprintf("select %s from orders o where 1=1 ", selectParam)

	var args []interface{}
	var condition string

	//检索条件
	//1.0 根据状态
	//1.1 根据状态 --> 在根据状态，还要考虑 搜索来源，区别app端搜索和seller端搜索
	if order.OrderStatus != -1 {

		if order.OrderStatus == 80 {
			//1.1.1 如果是app端并且搜索全部的订单，那么显示该用户所有状态的订单
			if order.SearchType != 0 {
				//1.1.3 如果是seller端并且搜索全部的订单，那么显示不为 代付款 ：0 和 已关闭订单：16 状态的订单
				condition += "and (o.order_status <> 0 and o.order_status <> 16) "
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
		condition += fmt.Sprintf(" and extract(epoch from o.pay_at)::integer between $%d and $%d", len(args), len(args)+1)
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

	condition += " order by o.update_at desc"
	if order.Page <= 0 {
		order.Page = 1
	}
	if order.Size <= 0 {
		order.Size = 10
	}
	condition += fmt.Sprintf(" OFFSET %d LIMIT %d ", (order.Page-1)*order.Size, order.Size)
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
		err = rows.Scan(&next.Id, &next.OrderStatus, &next.TotalFee, &next.Freight, &next.GoodsFee, &next.WithdrawalFee, &next.UserId, &next.Mobile, &next.Name, &next.Address, &next.Remark, &next.StoreId, &next.SchoolId, &next.TradeNo, &next.PayChannel, &next.OrderAt, &pay_at, &deliver_at, &print_at, &complete_at, &next.PrintStaffId, &next.DeliverStaffId, &next.AfterSaleStaffId, &after_sale_apply_at, &after_sale_end_at, &next.AfterSaleStatus, &next.AfterSaleTradeNo, &next.RefundFee, &next.GrouponId,
			&next.UpdateAt)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
		orderitems, err := GetOrderItems(next)
		if err != nil {
			misc.LogErr(err)
			return nil, err
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

	}

	return
}

//根据订单获取订单项集合
func GetOrderItems(order *pb.Order) (orderitems []*pb.OrderItem, err error) {
	query := "select oi.id,g.id,oi.type,oi.amount,oi.price,b.title,b.isbn,b.image from orders_item oi join goods g on oi.goods_id=g.id join books b on g.book_id=b.id where orders_id='%s'"
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
		err = rows.Scan(&orderItem.Id, &orderItem.GoodsId, &orderItem.Type, &orderItem.Amount, &orderItem.Price, &orderItem.BookTitle, &orderItem.BookIsbn, &orderItem.BookImage)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
	}
	return
}

//更改时间
func UpdateOrder(order *pb.Order) error {
	query := "update orders set update_at=now()"

	var args []interface{}
	var condition string
	if order.OrderStatus != 0 {
		args = append(args, order.OrderStatus)
		condition += fmt.Sprintf(",order_status=$%d", len(args))
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

//获取订单的状态
func GetOrderBaseInfo(order *pb.Order) error {
	next := order
	//需要的项
	var pay_at, deliver_at, print_at, complete_at, after_sale_apply_at, after_sale_end_at sql.NullString

	selectParam := "o.id,o.order_status,o.total_fee,o.freight,o.goods_fee,o.withdrawal_fee,o.user_id,o.mobile,o.name,o.address,o.remark,o.store_id,o.school_id,o.trade_no,o.pay_channel,extract(epoch from o.order_at)::integer,extract(epoch from o.pay_at)::integer,extract(epoch from o.deliver_at)::integer,extract(epoch from o.print_at)::integer,extract(epoch from o.complete_at)::integer,o.print_staff_id,o.deliver_staff_id,o.after_sale_staff_id,extract(epoch from o.after_sale_apply_at)::integer,extract(epoch from o.after_sale_end_at)::integer,o.after_sale_status,o.after_sale_trad_no,o.refund_fee,o.groupon_id,extract(epoch from o.update_at)::integer"

	query := fmt.Sprintf("select %s from orders o where id=$1 ", selectParam)

	err := DB.QueryRow(query).Scan(&next.Id, &next.OrderStatus, &next.TotalFee, &next.Freight, &next.GoodsFee, &next.WithdrawalFee, &next.UserId, &next.Mobile, &next.Name, &next.Address, &next.Remark, &next.StoreId, &next.SchoolId, &next.TradeNo, &next.PayChannel, &next.OrderAt, &pay_at, &deliver_at, &print_at, &complete_at, &next.PrintStaffId, &next.DeliverStaffId, &next.AfterSaleStaffId, &after_sale_apply_at, &after_sale_end_at, &next.AfterSaleStatus, &next.AfterSaleTradeNo, &next.RefundFee, &next.GrouponId,
		&next.UpdateAt)
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

	return nil
}
