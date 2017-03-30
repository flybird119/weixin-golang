package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

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
		log.Warn("================sdfhgkljdsfhgkldsfhkgjh==========")
		log.Warn("没找到")

		log.Error(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Warn(pb.LoginResp{Code: "00000", Message: "notFund"})
		log.Warn(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Debug(";;;;;;;;;;;;;;;")
		return &pb.LoginResp{Code: "00000", Message: "notFund"}, nil
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
	/**
	*====================================
	*		token 构建
	*====================================
	 */
	userinfo := &pb.UserInfo{Id: id, Username: in.Username, Mobile: in.Mobile, Token: "tokenya"}

	return &pb.RegisterResp{Code: "00000", Message: "ok", Data: userinfo}, nil
}
