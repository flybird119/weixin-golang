package service

import (
	"errors"

	"github.com/garyburd/redigo/redis"
	baseDb "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/misc/token"
	orderDb "github.com/goushuyun/weixin-golang/order/db"
	sellerDb "github.com/goushuyun/weixin-golang/seller/db"
	"github.com/wothing/log"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/seller/role"
	"github.com/goushuyun/weixin-golang/store/db"
)

//StoreServiceServer 店铺server
type StoreServiceServer struct{}

//AddStore 增加店铺
func (s *StoreServiceServer) AddStore(ctx context.Context, in *pb.Store) (*pb.AddStoreResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddStore", "%#v", in))

	//添加店铺
	err := db.AddStore(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//添加店铺和商家的关联
	err = db.AddStoreSellerMap(in, role.InterAdmin)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	in.AdminMobile = in.Seller.Mobile

	/**
	*================
	*	初始化操作
	* 	1: account 初始化 InitAccount
	*	2: 轮播图初始化
	*================
	 */

	account := &pb.Account{StoreId: in.Id, Type: 1}
	misc.CallRPC(ctx, "bc_account", "InitAccount", account)

	circular := &pb.Circular{StoreId: in.Id}
	misc.CallRPC(ctx, "bc_circular", "CircularInit", circular)

	/**
	*================
	*	记录日志
	*================
	 */
	return &pb.AddStoreResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UpdateStore 增加店铺
func (s *StoreServiceServer) UpdateStore(ctx context.Context, in *pb.Store) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateStore", "%#v", in))

	err := db.UpdateStore(in)

	/**
	*================
	*	记录日志
	*================
	 */

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//AddRealStore 增加实体店
func (s *StoreServiceServer) AddRealStore(ctx context.Context, in *pb.RealStore) (*pb.AddRealStoreResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddRealStore", "%#v", in))

	err := db.AddRealStore(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	/**
	*================
	*	记录日志
	*================
	 */
	return &pb.AddRealStoreResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UpdateRealStore 增加实体店
func (s *StoreServiceServer) UpdateRealStore(ctx context.Context, in *pb.RealStore) (*pb.NormalResp, error) {
	//记录踪迹
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateRealStore", "%#v", in))
	err := db.UpdateRealStore(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	/**
	*================
	*	记录日志
	*================
	 */
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//StoreInfo 查看店铺的信息
func (s *StoreServiceServer) StoreInfo(ctx context.Context, in *pb.Store) (*pb.AddStoreResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddStore", "%#v", in))
	err := db.GetStoreInfo(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.AddStoreResp{Code: "00000", Message: "ok", Data: in}, nil
}

//EnterStore 进入店铺
func (s *StoreServiceServer) EnterStore(ctx context.Context, in *pb.Store) (*pb.AddStoreResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "EnterStore", "%#v", in))

	err := db.GetStoreInfo(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	role, err := db.GetSellerStoreRole(in.Seller.Id, in.Id)
	if err != nil {
		log.Debug(err)
	}
	//重新牵token
	tokenStr := token.SignSellerToken(token.InterToken, in.Seller.Id, in.Seller.Mobile, in.Id, role)
	return &pb.AddStoreResp{Code: "00000", Message: "ok", Data: in, Token: tokenStr}, nil
}

//ChangeStoreLogo 修改店铺头像
func (s *StoreServiceServer) ChangeStoreLogo(ctx context.Context, in *pb.Store) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ChangeStoreLogo", "%#v", in))
	err := db.ChangeStoreLogo(in.Logo, in.Id)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//RealStores 获取实体店列表
func (s *StoreServiceServer) RealStores(ctx context.Context, in *pb.Store) (*pb.RealStoresResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddStore", "%#v", in))

	shops, err := db.GetStoreShops(in.Id)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.RealStoresResp{Code: "00000", Message: "ok", Data: shops}, nil
}

//CheckCode 校验验证码
func (s *StoreServiceServer) CheckCode(ctx context.Context, in *pb.RegisterModel) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CheckCode", "%#v", in))
	//获取redis的连接
	conn := baseDb.GetRedisConn()
	//检验验证码是否正确
	code, err := redis.String(conn.Do("get", "sellerUpdate:"+in.Mobile))
	if err == redis.ErrNil || code != in.MessageCode {
		log.Debugf("验证码错误：%s:%s", code, in.MessageCode)
		return &pb.NormalResp{Code: "00000", Message: "codeErr"}, nil
	}
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	conn.Close()

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//TransferStore 转让店铺
func (s *StoreServiceServer) TransferStore(ctx context.Context, in *pb.TransferStoreReq) (*pb.AddStoreResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "TransferStore", "%#v", in))
	//获取redis的连接
	conn := baseDb.GetRedisConn()
	defer conn.Close()
	//检验验证码是否正确
	code, err := redis.String(conn.Do("get", "sellerUpdate:"+in.Mobile))
	if err == redis.ErrNil || code != in.MessageCode {
		log.Debugf("验证码错误：%s:%s", code, in.MessageCode)
		return &pb.AddStoreResp{Code: "00000", Message: "codeErr"}, nil
	}
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	//转让店铺的核心操作
	//首先根据手机号获取用户的id
	sellerId, err := sellerDb.GetSellerByMobile(in.Mobile)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	err = db.TransferStore(sellerId, in.Store.Id)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.AddStoreResp{Code: "00000", Message: "ok"}, nil
}

//DelRealStore 删除实体店
func (s *StoreServiceServer) DelRealStore(ctx context.Context, in *pb.RealStore) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "TransferStore", "%#v", in))

	err := db.DelRealStore(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//记录日志
	log.Debugf("DelRealStore realStoreId:%s Operater Id:%s", in.Id, in.Seller.Id)
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

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

//更新提现账号
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

//更新提现账号
func (s *StoreServiceServer) StoreHistoryStateOrderNum(ctx context.Context, in *pb.StoreHistoryStateOrderNumModel) (*pb.StoreHistoryStateOrderNumResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StoreHistoryStateOrderNum", "%#v", in))

	err := orderDb.StoreCenterNecessaryOrderCount(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.StoreHistoryStateOrderNumResp{Code: "00000", Message: "ok", Data: in}, nil
}
