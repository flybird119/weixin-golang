package main

import (
	"github.com/goushuyun/weixin-golang/address/service"
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_address"
	port    = 8866
)

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterAddressServiceServer(s, &service.AddressServiceServer{})
	s.Serve(m.CreateListener())
}
