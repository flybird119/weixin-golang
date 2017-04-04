package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/misc"

	"github.com/wothing/log"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/seller/role"
	"github.com/goushuyun/weixin-golang/store/db"
)

//StoreServiceServer 店铺server
type StoreServiceServer struct{}

//AddStore 增加店铺
func (s *StoreServiceServer) AddStore(ctx context.Context, in *pb.Store) (*pb.AddStoreResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddStore", "%#v", in))

	//添加店铺
	err := db.AddStore(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//添加店铺和商家的关联
	err = db.AddStoreSellerMap(in, role.InterAdmin)
	if err != nil {
		log.Debug(err)
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
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateStore", "%#v", in))

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
