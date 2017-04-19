package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/goushuyun/weixin-golang/circular/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_circular"
	port    = 8862
)

var svcNames = []string{}

func main() {
	m := db.NewMicro(svcName, port)
	m.ReferServices(svcNames...)
	m.RegisterPG()
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterCircularServiceServer(s, &service.CircularServiceServer{})
	s.Serve(m.CreateListener())
}
