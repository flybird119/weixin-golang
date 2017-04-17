package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/weixin/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_weixin"
	port    = 8859
)

var svcNames = []string{}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.RegisterRedis()
	m.ReferServices(svcNames...)

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterWeixinServiceServer(s, &service.WeixinServer{})
	s.Serve(m.CreateListener())
}
