package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/user/db"
	"github.com/wothing/log"
)

type UserService struct {
}

func (s *UserService) SaveUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveUser", "%#v", req))

	// save user
	err := db.SaveUser(req)
	if err != nil {
		log.Error(err)
		errs.Wrap(errors.New(err.Error()))
	}

	return req, nil
}
