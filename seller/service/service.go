package service

import (
	"context"
	"fmt"

	"github.com/goushuyun/weixin-golang/pb"
)

type SellerServiceServer struct {
}

//UserPasswordLogin 登录
func (s *SellerServiceServer) UserPasswordLogin(ctx context.Context, in *pb.LoginReqModel) (*pb.LoginRspModel, error) {
	fmt.Println("111")
	return &pb.LoginRspModel{}, nil
}
