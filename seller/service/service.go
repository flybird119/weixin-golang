package service

import (
	"errors"

	"github.com/garyburd/redigo/redis"
	baseDb "github.com/goushuyun/weixin-golang/db"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/misc/token"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/seller/db"
	"github.com/goushuyun/weixin-golang/seller/role"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//SellerServiceServer server
type SellerServiceServer struct{}

//SellerLogin 登录
func (s *SellerServiceServer) SellerLogin(ctx context.Context, in *pb.LoginModel) (*pb.LoginResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SellerLogin", "%#v", in))

	sellerInfo, err := db.CheckSellerExists(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	if sellerInfo == nil {
		return &pb.LoginResp{Code: "00000", Message: "notFound"}, nil
	}
	/**
	*====================================
	*		token 构建
	*====================================
	 */
	tokenStr := token.SignSellerToken(token.InterToken, sellerInfo.Id, sellerInfo.Mobile, "", int64(role.AppNormalUser))
	sellerInfo.Token = tokenStr
	return &pb.LoginResp{Code: "00000", Message: "ok", Data: sellerInfo}, nil
}

//SellerRegister 商家用户注册
func (s *SellerServiceServer) SellerRegister(ctx context.Context, in *pb.RegisterModel) (*pb.RegisterResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SellerRegister", "%#v", in))

	/**
	*====================================
	*		检验手机验证码
	*====================================
	 */
	conn := baseDb.GetRedisConn()
	code, err := redis.String(conn.Do("get", "sellerRegister:"+in.Mobile))
	if err == redis.ErrNil || code != in.MessageCode {
		log.Debugf("验证码错误：%s:%s", code, in.MessageCode)
		return &pb.RegisterResp{Code: "00000", Message: "codeErr"}, nil
	}
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	conn.Close()

	id, err := db.SellerRegister(in)

	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//如果id为空，那么存在有相同注册的手机号，不再支持相同手机号的注册
	if id == "" {
		log.Debugf("存在手机号%s，不支持相同手机号的注册", in.Mobile)
		return &pb.RegisterResp{Code: "00000", Message: "exist"}, nil
	}
	/**
	*====================================
	*		token 构建
	*====================================
	 */
	tokenStr := token.SignSellerToken(token.InterToken, id, in.Mobile, "", int64(role.AppNormalUser))

	sellerInfo := &pb.SellerInfo{Id: id, Username: in.Username, Mobile: in.Mobile, Token: tokenStr}

	return &pb.RegisterResp{Code: "00000", Message: "ok", Data: sellerInfo}, nil
}

//UpdatePasswordAndLogin 修改登录密码
func (s *SellerServiceServer) UpdatePasswordAndLogin(ctx context.Context, in *pb.RegisterModel) (*pb.RegisterResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdatePasswordAndLogin", "%#v", in))

	/**
	*====================================
	*		检验手机验证码
	*====================================
	 */
	conn := baseDb.GetRedisConn()
	code, err := redis.String(conn.Do("get", "sellerUpdate:"+in.Mobile))
	if err == redis.ErrNil || code != in.MessageCode {
		log.Debugf("验证码错误：%s:%s", code, in.MessageCode)
		return &pb.RegisterResp{Code: "00000", Message: "codeErr"}, nil
	}
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	conn.Close()
	log.Debugf("===================,%s", 1)
	seller, err := db.UpdatePassword(in)

	if err != nil {
		log.Debugf("===================,%s", err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	log.Debugf("===================,%s", seller)
	/**
	*====================================
	*		token 构建
	*====================================
	 */
	tokenStr := token.SignSellerToken(token.InterToken, seller.Id, in.Mobile, "", int64(role.AppNormalUser))

	sellerInfo := &pb.SellerInfo{Id: seller.Id, Username: seller.Username, Mobile: in.Mobile, Token: tokenStr}

	return &pb.RegisterResp{Code: "00000", Message: "ok", Data: sellerInfo}, nil
}

//CheckMobileExist 检验手机号是否注册过
func (s *SellerServiceServer) CheckMobileExist(ctx context.Context, in *pb.CheckMobileReq) (*pb.CheckMobileRsp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CheckMobileExist", "%#v", in))

	isExist := db.CheckMobileExist(in.Mobile)
	if isExist {
		return &pb.CheckMobileRsp{Code: "00000", Message: "exist"}, nil
	}
	return &pb.CheckMobileRsp{Code: "00000", Message: "ok"}, nil
}

//GetTelCode 获取验证码
func (s *SellerServiceServer) GetTelCode(ctx context.Context, in *pb.CheckMobileReq) (*pb.CheckMobileRsp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetTelCode", "%#v", in))

	//检查手机号是否存在
	isExist := db.CheckMobileExist(in.Mobile)
	if isExist {
		return &pb.CheckMobileRsp{Code: "00000", Message: "exist"}, nil
	}
	code := misc.GenCheckCode(4, misc.KC_RAND_KIND_NUM)
	log.Debugf("sms code:%s", code)
	//rpc调用短信接口
	_, err := misc.CallRPC(ctx, "bc_sms", "SendSMS", &pb.SMSReq{Type: pb.SMSType_CommonCheckCode, Mobile: in.Mobile, Message: []string{code}})
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//redis 存放验证码
	conn := baseDb.GetRedisConn()
	_, err = conn.Do("set", "sellerRegister:"+in.Mobile, code)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	_, err = conn.Do("expire", "sellerRegister:"+in.Mobile, 300)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	conn.Close()
	return &pb.CheckMobileRsp{Code: "00000", Message: "ok"}, nil
}

//GetUpdateTelCode 获取修改密码的短信
func (s *SellerServiceServer) GetUpdateTelCode(ctx context.Context, in *pb.CheckMobileReq) (*pb.CheckMobileRsp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetUpdateTelCode", "%#v", in))

	//检查手机号是否存在
	isExist := db.CheckMobileExist(in.Mobile)
	if !isExist {
		return &pb.CheckMobileRsp{Code: "00000", Message: "needRegister"}, nil
	}
	code := misc.GenCheckCode(4, misc.KC_RAND_KIND_NUM)
	log.Debugf("sms code:%s", code)
	//rpc调用短信接口
	_, err := misc.CallRPC(ctx, "bc_sms", "SendSMS", &pb.SMSReq{Type: pb.SMSType_CommonCheckCode, Mobile: in.Mobile, Message: []string{code}})
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//redis 存放验证码
	conn := baseDb.GetRedisConn()
	_, err = conn.Do("set", "sellerUpdate:"+in.Mobile, code)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	_, err = conn.Do("expire", "sellerUpdate:"+in.Mobile, 300)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	conn.Close()
	return &pb.CheckMobileRsp{Code: "00000", Message: "ok"}, nil
}

//SelfStores 自己管理的店铺
func (s *SellerServiceServer) SelfStores(ctx context.Context, in *pb.SellerInfo) (*pb.SelfStoresResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SellerRegister", "%#v", in))

	store, err := db.GetStoresBySeller(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SelfStoresResp{Code: "00000", Message: "ok", Data: store}, nil
}
