package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/payment/service"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_payment"
	port    = 8865
)

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterPaymentServiceServer(s, &service.PaymentService{})
	s.Serve(m.CreateListener())
}
