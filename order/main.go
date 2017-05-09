package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/order/service"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/robfig/cron"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_order"
	port    = 8863
)

var svcNames = []string{
	"bc_cart",
	"bc_account",
	"bc_payment",
}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterOrderServiceServer(s, &service.OrderServiceServer{})

	//注册时间轮询
	c := cron.New()
	service.RegisterOrderPolling(c)
	c.Start()
	defer c.Stop()
	s.Serve(m.CreateListener())
}
