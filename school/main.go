package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/school/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_school"
	port    = 8853
)

var svcNames = []string{}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)
	// 注册redis
	db.InitRedis(svcName)
	defer db.CloseRedis()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterSchoolServiceServer(s, &service.SchoolServiceServer{})

	s.Serve(m.CreateListener())
}
