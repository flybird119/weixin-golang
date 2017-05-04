package service

import (
	"errors"
	"fmt"
	"strings"

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

//支付完成后，账户处理
func (s *AccountServiceServer) PayOverOrderAccountHandle(ctx context.Context, in *pb.Order) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "PayOverOrderAccountHandle", "%#v", in))

	hasSellerAccountItem, err := db.HasExistAcoount(&pb.AccountItem{StoreId: in.StoreId, OrderId: in.Id, ItemType: 4})
	if err != nil {
		go misc.LogErrOrder(in, "支付完成-检查记录是否存在发生错误,修改商户可提现余额时发生错误,影响商户支入支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if hasSellerAccountItem {
		go misc.LogErrOrder(in, "支付完成-需要检查订单", errors.New("修改失败，改数据已经在账单项中存在"))
		return nil, errs.Wrap(errors.New("修改失败，改数据已经在账单项中存在"))
	}
	//1.0 计算手续费
	serviceFee := in.TotalFee - in.WithdrawalFee
	//2.0 修改商户的可提现余额
	sellerAccount := &pb.Account{StoreId: in.StoreId, UnsettledBalance: in.WithdrawalFee}

	if sellerAccount.StoreId == "" {
		go misc.LogErrOrder(in, "支付完成-商户id为空,修改商户可提现余额时发生错误,影响商户支入支出记录和管理员余额修改异常及管理员入账记录", errors.New("store id is nil"))
		return nil, errs.Wrap(errors.New("store id is nil"))
	}
	err = db.ChangeAccountWithdrawalFee(sellerAccount)
	if err != nil {
		go misc.LogErrOrder(in, "支付完成-修改商户可提现余额时发生错误,影响商户支入支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//2.1 记录增加记录
	sellerRemark := "由%s支付，已经入待结算金额"
	if strings.Contains(in.PayChannel, "alipay") {
		sellerRemark = fmt.Sprintf(sellerRemark, "支付宝")
	} else if strings.Contains(in.PayChannel, "wx") {
		sellerRemark = fmt.Sprintf(sellerRemark, "微信")
	} else {
		sellerRemark = fmt.Sprintf(sellerRemark, "未知渠道")
	}
	sellerAccountItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 4, Remark: sellerRemark, ItemFee: in.TotalFee, AccountBalance: sellerAccount.UnsettledBalance + serviceFee}
	err = db.AddAccountItem(sellerAccountItem)
	if err != nil {
		go misc.LogErrAccount(sellerAccountItem, "支付完成-影响商户支入记录、商户手续费支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//2.2 服务费记录
	sellerServiceAccountItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 2, Remark: "在待结算金额中扣除", ItemFee: -serviceFee, AccountBalance: sellerAccount.UnsettledBalance}
	err = db.AddAccountItem(sellerServiceAccountItem)
	if err != nil {
		go misc.LogErrAccount(sellerServiceAccountItem, "支付完成-商户手续费支出记录和管理员余额修改异常及管理员入账记录", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//3 更改管理员账户的值
	adminAccount := &pb.Account{StoreId: "", UnsettledBalance: serviceFee}
	err = db.ChangeAccountWithdrawalFee(adminAccount)
	if err != nil {
		log.Error(err)
		go misc.LogErrOrder(in, "支付完成-管理员余额修改异常及管理员入账记录", err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//3.1 管理员相关记录
	adminAccountItem := &pb.AccountItem{UserType: 2, StoreId: "", OrderId: in.Id, ItemType: 33, Remark: "商家订单手续费", ItemFee: serviceFee, AccountBalance: adminAccount.UnsettledBalance}
	err = db.AddAccountItem(adminAccountItem)
	if err != nil {
		go misc.LogErrAccount(adminAccountItem, "支付完成-管理员记账异常", err)
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

//订单完成账户操作
// -> 订单完成包含两种两种情况
// -> > 1:订单经过正常流程完成（用户下单-用户支付-商家发货-用户确认订单或者系统确认订单)
// -> > 2:订单售后完成，并且在订单成功(订单状态在用户确认订单或者系统确认订单)之前
func (s *AccountServiceServer) OrderCompleteAccountHandle(ctx context.Context, in *pb.Order) (*pb.Void, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "OrderCompleteAccountHandle", "%#v", in))
	//首先查看相对应的记录
	hasSellerWithdrwalAccountItem, err := db.HasExistAcoount(&pb.AccountItem{StoreId: in.StoreId, OrderId: in.Id, ItemType: 1})
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-检查记录是否存在发生错误,影响商户可提现和待结算的转换问题", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	hasSellerBalanceAccountItem, err := db.HasExistAcoount(&pb.AccountItem{StoreId: in.StoreId, OrderId: in.Id, ItemType: 17})
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-检查记录是否存在发生错误 ,影响商户可提现和待结算的转换问题", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if hasSellerWithdrwalAccountItem || hasSellerBalanceAccountItem {
		go misc.LogErrOrder(in, "订单完成-需要检查订单，影响商户可提现和待结算的转换问题", errors.New("修改失败，改数据已经在账单项中存在"))
		return nil, errs.Wrap(errors.New("修改失败，改数据已经在账单项中存在"))
	}
	//1 需要更改两个值 1 ：待体现金额  2 可提现金额
	//2 需要记录两条记录 1 ：待体现金额资金流向记录  2 ： 可提现金额资金流入记录
	//1.1	待体现金额

	sellerAccountWithdrawal := &pb.Account{StoreId: in.StoreId, UnsettledBalance: -in.WithdrawalFee}
	err = db.ChangeAccountWithdrawalFee(sellerAccountWithdrawal)
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-待体现转化成可提现发生错误 ,影响商户可提现和待结算的转换", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	sellerAccountWithdrawalItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 1, Remark: "已由待结算金额转换为可提现金额", ItemFee: -in.WithdrawalFee, AccountBalance: sellerAccountWithdrawal.UnsettledBalance}

	//1.2 	待体现金额资金流向记录
	err = db.AddAccountItem(sellerAccountWithdrawalItem)
	if err != nil {
		go misc.LogErrAccount(sellerAccountWithdrawalItem, "订单完成-增加待结算转化可结算时发生错误 ,影响下一步更改可体现金额和增加资金流向记录的操作", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	sellerAccountBalance := &pb.Account{StoreId: in.StoreId, Balance: in.WithdrawalFee}
	//2.1：	待体现金额资金流向记录
	err = db.ChangAccountBalance(sellerAccountBalance)
	if err != nil {
		go misc.LogErrOrder(in, "订单完成-更改可结算金额操作时发生错误 ,影响商户可提现和待结算的转换", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//2.2	可提现金额资金流入记录
	sellerAccountBalanceItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 17, Remark: "由待结算金额转换为可提现金额", ItemFee: in.WithdrawalFee, AccountBalance: sellerAccountBalance.Balance}

	err = db.AddAccountItem(sellerAccountBalanceItem)
	if err != nil {
		go misc.LogErrAccount(sellerAccountWithdrawalItem, "订单完成-记录可结算操作时发生错误 ,影响操作日志", err)
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.Void{}, nil
}

//处理售后订单
func (s *AccountServiceServer) HandleAfterSaleOrder(ctx context.Context, in *pb.Order) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "PayOverOrderAccountHandle", "%#v", in))
	sellerAccountBalance := &pb.Account{StoreId: in.StoreId, Balance: -in.WithdrawalFee}
	//2.1：	待体现金额资金流向记录
	err := db.ChangAccountBalance(sellerAccountBalance)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//可提现售后24
	sellerAccountBalanceItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, OrderId: in.Id, ItemType: 24, Remark: "已从可提现金额扣除", ItemFee: -in.WithdrawalFee, AccountBalance: sellerAccountBalance.Balance}
	err = db.AddAccountItem(sellerAccountBalanceItem)
	if err != nil {
		log.Error(err)
		go misc.LogErrAccount(sellerAccountBalanceItem, "订单售后-增加记录失败", err)
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
