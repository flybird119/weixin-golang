package service

import (
	"database/sql"
	"errors"

	"github.com/garyburd/redigo/redis"
	accountDb "github.com/goushuyun/weixin-golang/account/db"
	baseDb "github.com/goushuyun/weixin-golang/db"
	sellerDb "github.com/goushuyun/weixin-golang/seller/db"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/wothing/log"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/store/db"
)

//获取提现账号管理的手机验证码
func (s *StoreServiceServer) GetCardOperSmsCode(ctx context.Context, in *pb.SmsCardSubmitModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetCardOperSmsCode", "%#v", in))
	store := &pb.Store{Id: in.StoreId}
	err := db.GetStoreInfo(store)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	in.Mobile = store.AdminMobile
	code := misc.GenCheckCode(4, misc.KC_RAND_KIND_NUM)

	log.Debugf("sms code:%s", code)
	//rpc调用短信接口
	_, err = misc.CallRPC(ctx, "bc_sms", "SendSMS", &pb.SMSReq{Type: pb.SMSType_CommonCheckCode, Mobile: store.AdminMobile, Message: []string{code}})
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//redis 存放验证码
	conn := baseDb.GetRedisConn()
	defer conn.Close()
	_, err = conn.Do("set", "storeCardOpe:"+in.Mobile, code)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	_, err = conn.Do("expire", "storeCardOpe:"+in.Mobile, 300)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//保存提现账号
func (s *StoreServiceServer) SaveStoreWithdrawCard(ctx context.Context, in *pb.StoreWithdrawCard) (*pb.StoreWithdrawCardOpeResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveStoreWithdrawCard", "%#v", in))
	store := &pb.Store{Id: in.StoreId}
	err := db.GetStoreInfo(store)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	in.Mobile = store.AdminMobile
	//首先检查code
	conn := baseDb.GetRedisConn()
	defer conn.Close()
	//检验验证码是否正确
	code, err := redis.String(conn.Do("get", "storeCardOpe:"+in.Mobile))
	if err == redis.ErrNil || code != in.Code {
		log.Debugf("验证码错误：%s:%s", code, in.Code)
		return &pb.StoreWithdrawCardOpeResp{Code: "00000", Message: "codeErr"}, nil
	}
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	err = db.SaveWithdrawCard(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.StoreWithdrawCardOpeResp{Code: "00000", Message: "ok", Data: in}, nil
}

//更新提现账号
func (s *StoreServiceServer) UpdateStoreWithdrawCard(ctx context.Context, in *pb.StoreWithdrawCard) (*pb.StoreWithdrawCardOpeResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateStoreWithdrawCard", "%#v", in))
	store := &pb.Store{Id: in.StoreId}
	err := db.GetStoreInfo(store)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	in.Mobile = store.AdminMobile
	//首先检查code
	conn := baseDb.GetRedisConn()
	defer conn.Close()
	//检验验证码是否正确
	code, err := redis.String(conn.Do("get", "storeCardOpe:"+in.Mobile))
	if err == redis.ErrNil || code != in.Code {
		log.Debugf("验证码错误：%s:%s", code, in.Code)
		return &pb.StoreWithdrawCardOpeResp{Code: "00000", Message: "codeErr"}, nil
	}
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	err = db.UpdateWithdrawCard(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.StoreWithdrawCardOpeResp{Code: "00000", Message: "ok", Data: in}, nil
}

//根据店铺获取提现账号
func (s *StoreServiceServer) GetWithdrawCardInfoByStore(ctx context.Context, in *pb.StoreWithdrawCard) (*pb.StoreWithdrawCardOpeResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetWithdrawCardInfoByStore", "%#v", in))

	err := db.GetWithdrawCardInfoByStore(in)

	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.StoreWithdrawCardOpeResp{Code: "00000", Message: "ok", Data: in}, nil
}

//提现申请
func (s *StoreServiceServer) WithdrawApply(ctx context.Context, in *pb.StoreWithdrawalsModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "WithdrawApply", "%#v", in))
	//开启事务
	tx, err := baseDb.DB.Begin()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//获取提现卡的信息
	card := &pb.StoreWithdrawCard{Id: in.WithdrawCardId}
	err = db.GetWithdrawCardInfoById(card)
	in.CardType = card.Type
	in.CardNo = card.CardNo
	in.CardName = card.CardName
	in.Username = card.Username

	defer tx.Rollback()

	seller, err := sellerDb.GetSellerById(in.StaffId)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	in.ApplyPhone = seller.Mobile
	//扣除可提现金额
	err = handleWithdrawApply(tx, in)
	if err != nil && err.Error() == "sellerNoMoney" {
		return &pb.NormalResp{Code: "00000", Message: "sellerNoMoney"}, nil
	}
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//申请记录
	err = db.SaveWithdrawApply(tx, in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	tx.Commit()
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil

}

//处理售后订单
func handleWithdrawApply(tx *sql.Tx, in *pb.StoreWithdrawalsModel) error {
	sellerAccountBalance := &pb.Account{StoreId: in.StoreId, Balance: -in.WithdrawFee}
	//2.1：	待体现金额资金流向记录
	err := accountDb.ChangAccountBalanceWithTx(tx, sellerAccountBalance)
	if err != nil {
		log.Error(err)
		return err
	}

	cardNo := in.CardNo
	if len(cardNo) < 4 {
		return errors.New("Card No InCorrect")
	}
	subCardNo := misc.SubString(cardNo, len(cardNo)-4, 4)
	//提现记录
	sellerAccountBalanceItem := &pb.AccountItem{UserType: 1, StoreId: in.StoreId, ItemType: 20, Remark: "提现到银行卡,尾号:" + subCardNo, ItemFee: -in.WithdrawFee, AccountBalance: sellerAccountBalance.Balance}
	err = accountDb.AddAccountItemWithTx(tx, sellerAccountBalanceItem)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
