package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	orderDB "github.com/goushuyun/weixin-golang/order/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type OrderServiceServer struct{}

// 提交订单
func (s *OrderServiceServer) OrderSubmit(ctx context.Context, in *pb.OrderSubmitModel) (*pb.OrderSubmitResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CartList", "%#v", in))

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
	defer log.TraceOut(log.TraceIn(tid, "CartList", "%#v", in))
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

	return &pb.NormalResp{}, nil
}

// 订单配送
func (s *OrderServiceServer) DistributeOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	return &pb.NormalResp{}, nil
}

// 确认订单（微信端）
func (s *OrderServiceServer) ConfirmOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	return &pb.NormalResp{}, nil
}

// 获取订单详情
func (s *OrderServiceServer) OrderDetail(ctx context.Context, in *pb.Order) (*pb.OrderDetailResp, error) {
	return &pb.OrderDetailResp{}, nil
}

// 获取订单列表 用户 云店铺 状态
func (s *OrderServiceServer) OrderList(ctx context.Context, in *pb.Order) (*pb.OrderListResp, error) {

	return &pb.OrderListResp{}, nil
}

//关闭订单
func (s *OrderServiceServer) CloseOrder(ctx context.Context, in *pb.Order) (*pb.Void, error) {

	return &pb.Void{}, nil
}
