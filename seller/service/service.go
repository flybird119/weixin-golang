package service

import (
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type SellerServiceServer struct{}

//UserPasswordLogin 登录
func (s *SellerServiceServer) UserPasswordLogin(ctx context.Context, in *pb.LoginReqModel) (*pb.LoginRspModel, error) {

	log.Debug(">>>>>>>>>akjsdhfklashdflkhl>>>>>>>>>>>")

	return &pb.LoginRspModel{}, nil
}
