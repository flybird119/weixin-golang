package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/location/service"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_location"
	port    = 8854
)

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterLocationServiceServer(s, &service.LocationServiceServer{})
	s.Serve(m.CreateListener())
}
