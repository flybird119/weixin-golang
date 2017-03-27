package main

import (
	"strings"

	"google.golang.org/grpc"

	"goushuyun/db"
	"goushuyun/mediastore/service"
	"goushuyun/pb"

	"github.com/wothing/worpc"
)

const svcName = "mediastore"
const port = 10017

func main() {
	m := db.NewMicro(svcName, port)

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	test := strings.ToLower(db.GetValue(svcName, "mode", "test")) != "live"
	pb.RegisterMediastoreServer(s, &service.MediastoreServer{Test: test})
	s.Serve(m.CreateListener())
}
