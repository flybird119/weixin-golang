package service

import (
<<<<<<< HEAD
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/misc"
=======
	"database/sql"
	"errors"

	accountDb "github.com/goushuyun/weixin-golang/account/db"
	"github.com/goushuyun/weixin-golang/errs"
>>>>>>> 4a9f0722db9e98903d030f85313d9f1ecb0c7ebf

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/store/db"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//充值申请
func (s *StoreServiceServer) RechargeApply(ctx context.Context, in *pb.RechargeModel) (*pb.RechargeOperationResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "RechargeApply", "%#v", in))
	err := db.RechargeApply(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.RechargeOperationResp{Code: "00000", Message: "ok", Data: in}, nil
}

//充值完成处理
func (s *StoreServiceServer) RechargeHandler(ctx context.Context, in *pb.RechargeModel) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "RechargeHandler", "%#v", in))

	searchRecharge := &pb.RechargeModel{Id: in.Id}
	err := db.GetRechargeById(searchRecharge)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//获取事务
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	defer tx.Rollback()
	//填写充值结果
	err = db.RechargeSuccessHandler(tx, in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//更改商户充值金额
	err = handleRechargeResult(tx, in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.Void{}, nil
}

//处理售后订单
func handleRechargeResult(tx *sql.Tx, in *pb.RechargeModel) error {
	sellerAccountBalance := &pb.Account{StoreId: in.StoreId, Balance: in.RechargeFee}
	//2.1：	待体现金额资金流向记录
	err := accountDb.ChangAccountBalanceWithTx(tx, sellerAccountBalance)
	if err != nil {
		log.Error(err)
		return err
	}

	//提现记录
	sellerAccountBalanceItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, ItemType: 18, Remark: "商家充值到可提现金额", ItemFee: in.RechargeFee, AccountBalance: sellerAccountBalance.Balance}
	err = accountDb.AddAccountItemWithTx(tx, sellerAccountBalanceItem)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
