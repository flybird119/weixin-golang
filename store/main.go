package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/store/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_store"
	port    = 8851
)

var svcNames = []string{
	"bc_sms",
	"bc_circular",
	"bc_account",
	"bc_mediastore",
}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)
	// 注册redis
	db.InitRedis(svcName)
	defer db.CloseRedis()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterStoreServiceServer(s, &service.StoreServiceServer{})

	s.Serve(m.CreateListener())
}
