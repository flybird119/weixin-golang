package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"

	"github.com/goushuyun/weixin-golang/groupon/service"
	"github.com/goushuyun/weixin-golang/pb"
)

const (
	svcName = "bc_groupon"
	port    = 8871
)

var svcNames = []string{}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterGrouponServiceServer(s, &service.GrouponServiceServer{})
	s.Serve(m.CreateListener())
}
