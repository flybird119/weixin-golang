package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/master/service"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_master"
	port    = 8869
)

var svcNames = []string{}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterMasterServiceServer(s, &service.MasterServiceServer{})
	s.Serve(m.CreateListener())
}
