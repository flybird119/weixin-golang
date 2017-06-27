package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"

	"github.com/goushuyun/weixin-golang/goods/service"
	"github.com/goushuyun/weixin-golang/pb"
)

const (
	svcName = "bc_goods"
	port    = 8856
)

var svcNames = []string{
	"bc_books",
	"bc_mediastore",
}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterGoodsServiceServer(s, &service.GoodsServiceServer{})
	s.Serve(m.CreateListener())
}
