package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/goushuyun/weixin-golang/cart/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_cart"
	port    = 8860
)

var svcNames = []string{
	"bc_goods",
}

func main() {
	m := db.NewMicro(svcName, port)
	m.ReferServices(svcNames...)
	m.RegisterPG()
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterCartServiceServer(s, &service.CartServiceServer{})
	s.Serve(m.CreateListener())
}
