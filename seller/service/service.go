package service

import (
	"github.com/goushuyun/weixin-golang/pb"
	"golang.org/x/net/context"
)

type SellerServiceServer struct{}

//UserPasswordLogin 登录
func (s *SellerServiceServer) SellerLogin(ctx context.Context, in *pb.LoginModel) (*pb.UserInfo, error) {

	return &pb.UserInfo{}, nil
}

//SellerRegister 商家用户注册
func (s *SellerServiceServer) SellerRegister(ctx context.Context, in *RegisterModel) (*UserInfo, error) {

	return &pb.UserInfo{}, nil
}
