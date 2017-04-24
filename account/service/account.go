package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/account/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type AccountServiceServer struct{}

//初始化account
func (s *AccountServiceServer) InitAccount(ctx context.Context, in *pb.Account) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "InitAccount", "%#v", in))
	err := db.InitAccount(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.Void{}, nil
}
func OrderAccountHandle(order *pb.Order) (*pb.Void, error) {

}

//增加什么AccountItem
func (s *AccountServiceServer) AddAccountItem(ctx context.Context, in *pb.AccountItem) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "InitAccount", "%#v", in))
	if item_type == 0 {

	}
	//misc.LogErrOrder(order, impact, err)
	return &pb.Void{}, nil
}
