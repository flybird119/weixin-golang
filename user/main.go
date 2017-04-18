package user

import (
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/user/service"
	"github.com/wothing/worpc"
	"google.golang.org/grpc"
)

const (
	svcName = "bc_user"
	port    = 8861
)

func main() {
	m := db.NewMicro(svcName, port)
	m.RegisterPG()

	s := grpc.NewServer(grpc.UnaryInterceptor(worpc.UnaryInterceptorChain(worpc.Recovery, worpc.Logging)))
	pb.RegisterUserServiceServer(s, &service.UserService{})
	s.Serve(m.CreateListener())
}
