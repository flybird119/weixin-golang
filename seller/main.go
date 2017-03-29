package main

import (
	"fmt"
	"net"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/seller/service"
	"github.com/wothing/log"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "seller"
	port    = 10014
)

func main() {

	db.InitPG(svcName)
	defer db.ClosePG()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", db.GetPort(port)))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Infof("starting to listen at : %d", db.GetPort(port))

	//registe admin to etcd
	err = db.RegisterService(svcName, port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterSellerServiceServer(s, &service.SellerServiceServer)
	s.Serve(lis)

}
