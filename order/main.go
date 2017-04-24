package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/order/service"
	"github.com/goushuyun/weixin-golang/pb"
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
}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterOrderServiceServer(s, &service.OrderServiceServer{})

	s.Serve(m.CreateListener())
}
