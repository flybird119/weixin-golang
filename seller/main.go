package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/seller/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_seller"
	port    = 8849
)

var svcNames = []string{
	"bc_sms",
	"bc_account",
}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)
	// 注册redis
	db.InitRedis(svcName)
	defer db.CloseRedis()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterSellerServiceServer(s, &service.SellerServiceServer{})

	s.Serve(m.CreateListener())
}
