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

//搜索罗列订单
