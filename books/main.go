package main

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"

	"github.com/goushuyun/weixin-golang/books/service"
	"github.com/goushuyun/weixin-golang/pb"
)

const (
	svcName = "bc_books"
	port    = 8855
)

var svcNames = []string{
	"bc_mediastore",
}

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()
	m.ReferServices(svcNames...)
	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))

	pb.RegisterBooksServiceServer(s, &service.BooksServer{})
	s.Serve(m.CreateListener())
}
