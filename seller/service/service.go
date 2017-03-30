package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/seller/db"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type SellerServiceServer struct{}

//UserPasswordLogin 登录
func (s *SellerServiceServer) SellerLogin(ctx context.Context, in *pb.LoginModel) (*pb.LoginResp, error) {

	userinfo, err := db.CheckSellerExists(in)

	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	/**
	*====================================
	*		token 构建
	*====================================
	 */
	if userinfo == nil {
		return &pb.LoginResp{Code: "00000", Message: "notFound"}, nil
	}
	log.Warn(*userinfo)
	userinfo.Token = "还没有token啦"

	return &pb.LoginResp{Code: "00000", Message: "ok", Data: userinfo}, nil
}

//SellerRegister 商家用户注册
func (s *SellerServiceServer) SellerRegister(ctx context.Context, in *pb.RegisterModel) (*pb.RegisterResp, error) {
	/**
	*====================================
	*		检验手机验证码
	*====================================
	 */
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
	userinfo := &pb.UserInfo{Id: id, Username: in.Username, Mobile: in.Mobile, Token: "tokenya"}

	return &pb.RegisterResp{Code: "00000", Message: "ok", Data: userinfo}, nil
}

//CheckMobileExist 检验手机号是否注册过
func (s *SellerServiceServer) CheckMobileExist(ctx context.Context, in *pb.CheckMobileReq) (*pb.CheckMobileRsp, error) {
	isExist := db.CheckMobileExist(in.Mobile)
	if isExist {
		return &pb.CheckMobileRsp{Code: "00000", Message: "exist"}, nil
	}
	return &pb.CheckMobileRsp{Code: "00000", Message: "ok"}, nil
}

//CheckMobileExist 检验手机号是否注册过
func (s *SellerServiceServer) GetTelCode(ctx context.Context, in *pb.CheckMobileReq) (*pb.CheckMobileRsp, error) {
	isExist := db.CheckMobileExist(in.Mobile)
	if isExist {
		return &pb.CheckMobileRsp{Code: "00000", Message: "exist"}, nil
	}
	code := misc.GenCheckCode(misc.KC_RAND_KIND_NUM, 4)

	_, err := misc.CallRPC(ctx, "sms", "SendSMS", &pb.SMSReq{Type: pb.SMSType_CommonCheckCode, Mobile: in.Mobile, Message: []string{code}})
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.CheckMobileRsp{Code: "00000", Message: "ok"}, nil
}
