package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/goushuyun/weixin-golang/errs"

	accountDB "github.com/goushuyun/weixin-golang/account/db"
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	orderDB "github.com/goushuyun/weixin-golang/order/db"
	storeDB "github.com/goushuyun/weixin-golang/store/db"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type OrderServiceServer struct{}

// 提交订单
func (s *OrderServiceServer) OrderSubmit(ctx context.Context, in *pb.OrderSubmitModel) (*pb.OrderSubmitResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "OrderSubmit", "%#v", in))

	//获取购物车
	req := &pb.Cart{Ids: in.CartIds, StoreId: in.StoreId, UserId: in.UserId}
	resp, err := misc.CallRPC(ctx, "bc_cart", "CartBaseList", req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	cartListResp, ok := resp.(*pb.CartListResp)
	if !ok {
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	carts := cartListResp.Data
	if len(carts) <= 0 {
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	tx, err := db.DB.Begin()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	defer misc.RollbackCtx(ctx, tx)
	//保存订单
	order, noStack, err := orderDB.OrderSubmit(tx, carts, in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//库存不足
	if noStack != "" {
		return &pb.OrderSubmitResp{Code: "00000", Message: "noStack", Data: order}, nil
	}
	//清理购物车
	misc.CallRPC(ctx, "bc_cart", "CartDel", req)
	tx.Commit()
	return &pb.OrderSubmitResp{Code: "00000", Message: "ok", Data: order}, nil
}

// 订单支付完成 -->待发货
func (s *OrderServiceServer) PaySuccess(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "PaySuccess", "%#v", in))
	//成功支付 准确记录值，如果其中一步发生错误,事务不会滚，
	//1 更改订单状态,填写支付方式和交易号， --异常，下面的事务不执行，写入操作异常
	isChanged, err := orderDB.PaySuccess(in)
	if err != nil {
		//--异常，下面的事务不执行，写入操作异常
		log.Error(err)
		misc.LogErrOrder(in, "更改订单支付状态时发生错误，影响商户支入支出记录和管理员余额修改异常及管理员入账记录", err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//1.2检查这个订单有没有被修改过
	if isChanged {
		log.Warn("order:%s hasChanged", in.Id)
		return &pb.NormalResp{Code: "00000", Message: "isChanged"}, nil
	}

	//2 修改商家账户和管理员账户 以及记录交易记录

	misc.CallRPC(ctx, "bc_account", "PayOverOrderAccountHandle", in)

	return &pb.NormalResp{}, nil
}

// 订单发货
func (s *OrderServiceServer) DeliverOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	//更改订单状态 订单状态 +2
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "DeliverOrder", "%#v", in))
	//首先检查发货前的订单状态
	searchOrder := &pb.Order{Id: in.Id}
	err := orderDB.GetOrderBaseInfo(searchOrder)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//开始检查订单的状态
	if searchOrder.OrderStatus != 1 {
		if err != nil {
			return nil, errs.Wrap(errors.New("order status error"))
		}
	}
	//填写操作人 并填写发送的时间和更改时间
	in.DeliverStaffId = in.SellerId
	in.DeliverAt = 1
	in.OrderStatus = 2
	in.UpdateAt = 1
	err = orderDB.UpdateOrder(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	//发送短信
	//发送短信 --- 获取店铺名称
	store := &pb.Store{Id: searchOrder.StoreId}
	err = storeDB.GetStoreInfo(store)
	if err == nil {
		//构建发送模版
		message := []string{store.Name, searchOrder.Id}
		_, err = misc.CallRPC(ctx, "bc_sms", "SendSMS", &pb.SMSReq{Type: pb.SMSType_Delivery, Mobile: searchOrder.Mobile, Message: message})
		if err != nil {
			log.Error(err)
		}

	} else {
		log.Error(err)
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

// 订单配送
func (s *OrderServiceServer) DistributeOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "DistributeOrder", "%#v", in))
	//填写操作人 并填写配送的时间并填写配送的时间和更改时间
	in.DistributeStaffId = in.SellerId
	in.DistributeAt = 1
	in.UpdateAt = 1
	err := orderDB.UpdateOrder(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

// 确认订单（微信端）
func (s *OrderServiceServer) ConfirmOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ConfirmOrder", "%#v", in))
	//1.0 首先要检验 订单的状态 未发货的订单不能点击成功
	searchOrder := &pb.Order{Id: in.Id}
	err := orderDB.GetOrderBaseInfo(searchOrder)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if searchOrder.OrderStatus != 3 || searchOrder.ConfirmAt != 0 {
		return nil, errs.Wrap(errors.New("order state error"))
	}
	//用户主动确认订单
	//订单成功——>修改订单状态 +4
	in.ConfirmAt = 1
	in.OrderStatus = 4
	in.CompleteAt = 1
	in.UpdateAt = 1
	err = orderDB.UpdateOrder(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	searchOrder.OrderStatus = 4
	//商家账户更改
	misc.CallRPC(ctx, "bc_account", "OrderCompleteAccountHandle", searchOrder)

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

// 申请售后（微信端）
func (s *OrderServiceServer) AfterSaleApply(ctx context.Context, in *pb.AfterSaleModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "PrintOrder", "%#v", in))
	//检查用户有没有资格申请售后
	serachOrder := &pb.Order{Id: in.OrderId}
	err := orderDB.GetOrderBaseInfo(serachOrder)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//1 如果订单状态是0，那么订单不支持售后
	if serachOrder.OrderStatus == 0 || serachOrder.OrderStatus == 8 {

		return nil, errs.Wrap(errors.New("order not support after-sales service"))
	}
	//2 完成订单14天之后不能申请售后
	now := time.Now()
	now = now.Add(14 * 24 * time.Hour)
	if serachOrder.CompleteAt > now.Unix() {
		return nil, errs.Wrap(errors.New("order not support after-sales service"))
	}
	//3 如果已经申请售后过了，就不能再次申请售后
	if serachOrder.AfterSaleStatus != 0 {
		return nil, errs.Wrap(errors.New("can not repeated apply after-sales service"))
	}
	//满足售后条件,那么检查必要字段
	//更改订单，修改after_sale_apply_at，after_sale_status，refund_fee,refund_fee,reason ,images
	err = orderDB.FillInAfterSaleDetail(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New("undefind err"))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

// 订单配送
func (s *OrderServiceServer) PrintOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	//打印时间
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "PrintOrder", "%#v", in))
	in.PrintAt = 1
	in.PrintStaffId = in.SellerId
	err := orderDB.UpdateOrder(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

// 获取订单详情
func (s *OrderServiceServer) OrderDetail(ctx context.Context, in *pb.Order) (*pb.OrderDetailResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "OrderDetail", "%#v", in))
	//首先获取订单的基本信息
	err := orderDB.GetOrderBaseInfo(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	orderitems, err := orderDB.GetOrderItems(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	staffs, err := orderDB.GetOrderStaffWork(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//获取售后详情
	afterSaleModel, err := orderDB.GetAfterSaleDetail(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if in.AfterSaleStaffId != "" {
		for i := 0; i < len(staffs); i++ {
			if staffs[i].StaffWork == "after_sale" {
				afterSaleModel.StaffId = staffs[i].StaffId
				afterSaleModel.StaffName = staffs[i].StaffName
			}
		}
	}
	orderDetail := &pb.OrderDetail{Order: in, Items: orderitems, AfterSaleDetail: afterSaleModel, Staffs: staffs}

	return &pb.OrderDetailResp{Code: "00000", Message: "ok", Data: orderDetail}, nil
}

//订单成功

// 获取订单列表 用户 云店铺 状态
func (s *OrderServiceServer) OrderList(ctx context.Context, in *pb.Order) (*pb.OrderListResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "OrderList", "%#v", in))
	details, err, totalCount := orderDB.FindOrders(in)
	if err != nil {
		log.Warn(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.OrderListResp{Code: "00000", Message: "ok", Data: details, TotalCount: totalCount}, nil
}

//关闭订单
func (s *OrderServiceServer) CloseOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CloseOrder", "%#v", in))
	//释放图书资源，更改修改过时间 更改订单状态
	//首先更改改订单的状态
	err := orderDB.CloseOrder(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//处理售后订单--未完成
func (s *OrderServiceServer) HandleAfterSaleOrder(ctx context.Context, in *pb.AfterSaleModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "HandleAfterSaleOrder", "%#v", in))
	order := &pb.Order{Id: in.OrderId}
	err := orderDB.GetOrderBaseInfo(order)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if order.AfterSaleStatus != 1 {
		return nil, errs.Wrap(errors.New("order not support after sale service"))
	}
	//order_id return_fee
	order.RefundFee = in.RefundFee
	tx, err := db.DB.Begin()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	defer tx.Rollback()
	log.Debug("hello")
	//检查退款金额 ，refund_fee == 0 ? 特殊处理：CallRPC
	err = handleAfterSaleOrder(tx, order)
	log.Debug("hello")
	//检查用户资金以及记录
	if err != nil && err.Error() == "sellerNoMoney" {
		return &pb.NormalResp{Code: "00000", Message: "可提现金额不足，请充值"}, nil
	} else if err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	log.Debug("ceshi320")
	//修改退款状态
	order.AfterSaleStaffId = in.StaffId
	err = orderDB.HandleAfterSaleOrder(tx, order)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	log.Debug("ceshi328")
	//查看退款金额是否为0
	if order.RefundFee != 0 {
		log.Debug("ceshi331")
		afterSaleModel, err := orderDB.GetAfterSaleDetail(order)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
		//如果退款金额不为0 ，Callrpc
		req := &pb.RefundReq{TradeNo: order.TradeNo, Amount: in.RefundFee, Reason: afterSaleModel.Reason}
		log.Debug("ceshi339")
		_, err = misc.CallRPC(ctx, "bc_payment", "Refund", req)
		log.Debug("ceshi340")
		if err != nil {
			log.Error(err)
			return &pb.NormalResp{Code: "00000", Message: err.Error()}, nil
		}

	} else {
		_, err := misc.CallRPC(ctx, "bc_account", "OrderCompleteAccountHandle", order)
		if err != nil && err.Error() != "code: 10000, message:修改失败，改数据已经在账单项中存在" {
			log.Error(err)
			return &pb.NormalResp{Code: "00000", Message: err.Error()}, nil
		}
	}
	tx.Commit()
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//售后订单处理结果
func (s *OrderServiceServer) AfterSaleOrderHandledResult(ctx context.Context, in *pb.AfterSaleModel) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AfterSaleOrderHandledResult", "%#v", in))
	order := &pb.Order{TradeNo: in.TradeNo}
	err := orderDB.GetOrderBaseInfoByTradeNo(order)
	if err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	if order.AfterSaleStatus <= 0 {
		return &pb.Void{}, nil
	}
	//检查是否处理过
	if order.AfterSaleStatus > 2 {
		return &pb.Void{}, nil
	}
	err = orderDB.AfterSaleResultOperation(in)
	if err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//检查是否满足体现条件
	if order.OrderStatus >= 17 && order.OrderStatus < 23 {
		//商家账户更改
		misc.CallRPC(ctx, "bc_account", "OrderCompleteAccountHandle", order)
	}
	return &pb.Void{}, nil
}

//处理售后订单
func (s *OrderServiceServer) UserCenterNecessaryOrderCount(ctx context.Context, in *pb.UserCenterOrderCount) (*pb.UserCenterOrderCount, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UserCenterNecessaryOrderCount", "%#v", in))
	err := orderDB.UserCenterNecessaryOrderCount(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return in, nil
}

// 导出发货订单
func (s *OrderServiceServer) ExportDeliveryOrderData(ctx context.Context, in *pb.Order) (*pb.OrderListResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ExportOrderData", "%#v", in))
	details, err := orderDB.ExportDeliveryOrderData(in)
	if err != nil {
		log.Warn(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.OrderListResp{Code: "00000", Message: "ok", Data: details}, nil
}

// 导出配货单
func (s *OrderServiceServer) ExportDistributeOrderData(ctx context.Context, in *pb.Order) (*pb.DistributeOrdersResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ExportDistributeOrderData", "%#v", in))
	models, err := orderDB.ExportDistributeOrderData(in)
	if err != nil {
		log.Warn(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.DistributeOrdersResp{Code: "00000", Message: "ok", Data: models}, nil
}

// 获导出订单
func (s *OrderServiceServer) OrderShareOperation(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "OrderShareOperation", "%#v", in))
	userId := in.UserId
	err := orderDB.GetOrderBaseInfo(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	orderitems, err := orderDB.GetOrderItems(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	for i := 0; i < len(orderitems); i++ {
		orderitem := orderitems[i]
		cart := &pb.Cart{UserId: userId, StoreId: in.StoreId, GoodsId: orderitem.GoodsId, Type: orderitem.Type, Amount: orderitem.Amount}
		misc.CallRPC(ctx, "bc_cart", "CartAdd", cart)
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//处理售后订单
func handleAfterSaleOrder(tx *sql.Tx, in *pb.Order) error {

	sellerAccountBalance := &pb.Account{StoreId: in.StoreId, Balance: -in.RefundFee}

	//如果退款0元，那么不记录
	if in.RefundFee == 0 {

		return nil
	}

	err := accountDB.ChangAccountBalanceWithTx(tx, sellerAccountBalance)
	if err != nil {
		log.Error(err)
		return err
	}
	//可提现售后24
	sellerAccountBalanceItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 24, Remark: "已从可提现金额扣除", ItemFee: -in.RefundFee, AccountBalance: sellerAccountBalance.Balance}
	err = accountDB.AddAccountItemWithTx(tx, sellerAccountBalanceItem)
	if err != nil {
		log.Error(err)
		go misc.LogErrAccount(sellerAccountBalanceItem, "订单售后-增加记录失败", err)
	}
	log.Debug("hello")
	return nil
}
