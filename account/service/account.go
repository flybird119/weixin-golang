package service

import (
	"errors"
	"fmt"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/account/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type AccountServiceServer struct{}

//初始化account
func (s *AccountServiceServer) InitAccount(ctx context.Context, in *pb.Account) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "InitAccount", "%#v", in))
	err := db.InitAccount(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.Void{}, nil
}
func (s *AccountServiceServer) PayOverOrderAccountHandle(ctx context.Context, in *pb.Order) (*pb.Void, error) {
	//1.0 计算手续费
	serviceFee := in.TotalFee - in.WithdrawalFee
	//2.0 修改商户的可提现余额
	sellerAccount := &pb.Account{StoreId: in.StoreId, UnsettledBalance: in.WithdrawalFee}

	if sellerAccount.StoreId == "" {
		go misc.LogErrOrder(in, "商户id为空,修改商户可提现余额时发生错误,影响商户支入支出记录和管理员余额修改异常及管理员入账记录", errors.New("store id is nil"))
		return nil, errs.Wrap(errors.New("store id is nil"))
	}
	err := db.ChangeAccountWithdrawalFee(sellerAccount)
	if err != nil {
		go misc.LogErrOrder(in, "修改商户可提现余额时发生错误,影响商户支入支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//2.1 记录增加记录
	sellerRemark := "由%s支付，已经入待结算金额"
	if in.PayChannel == "alipay" {
		sellerRemark = fmt.Sprintf(sellerRemark, "支付宝")
	} else {
		sellerRemark = fmt.Sprintf(sellerRemark, "微信")
	}
	sellerAccountItem := &pb.AccountItem{UserType: 0, StoreId: in.StoreId, OrderId: in.Id, ItemType: 1, Remark: sellerRemark, ItemFee: in.TotalFee, AccountBalance: sellerAccount.UnsettledBalance + serviceFee}
	err = db.AddAccountItem(sellerAccountItem)
	if err != nil {
		go misc.LogErrAccount(sellerAccountItem, "影响商户支入记录、商户手续费支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//2.2 服务费记录
	sellerServiceAccountItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 2, Remark: "在待结算金额中扣除", ItemFee: -serviceFee, AccountBalance: sellerAccount.UnsettledBalance}
	err = db.AddAccountItem(sellerServiceAccountItem)
	if err != nil {
		go misc.LogErrAccount(sellerServiceAccountItem, "商户手续费支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//3 更改管理员账户的值
	adminAccount := &pb.Account{StoreId: "", UnsettledBalance: serviceFee}
	err = db.ChangeAccountWithdrawalFee(adminAccount)
	if err != nil {
		log.Error(err)
		go misc.LogErrOrder(in, "管理员余额修改异常及管理员入账记录", err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//3.1 管理员相关记录
	adminAccountItem := &pb.AccountItem{UserType: 2, StoreId: "", OrderId: in.Id, ItemType: 33, Remark: "商家订单手续费", ItemFee: serviceFee, AccountBalance: adminAccount.UnsettledBalance}
	err = db.AddAccountItem(adminAccountItem)
	if err != nil {
		go misc.LogErrAccount(adminAccountItem, "管理员记账异常", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.Void{}, nil
}

//增加什么AccountItem
func (s *AccountServiceServer) AddAccountItem(ctx context.Context, in *pb.AccountItem) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "InitAccount", "%#v", in))
	//misc.LogErrOrder(order, impact, err)
	return &pb.Void{}, nil
}
