package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/statistic/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_statistic"
	port    = 8868
)

var svcNames = []string{}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterStatisticServiceServer(s, &service.StatisticServiceServer{})

	s.Serve(m.CreateListener())
}
