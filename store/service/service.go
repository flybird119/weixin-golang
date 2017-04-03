package service

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/store/db"
)

//StoreServiceServer 店铺server
type StoreServiceServer struct{}

//AddStore 增加店铺
func (s *StoreServiceServer) AddStore(ctx context.Context, in *pb.Store) (*pb.AddStoreResp, error) {
	err := db.AddStore(in)
	if err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	/**
	*================
	*	记录日志
	*================
	 */

	return &pb.AddStoreResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UpdateStore 增加店铺
func (s *StoreServiceServer) UpdateStore(ctx context.Context, in *pb.Store) (*pb.NormalResp, error) {
	err := db.UpdateStore(in)

	/**
	*================
	*	记录日志
	*================
	 */

	if err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
