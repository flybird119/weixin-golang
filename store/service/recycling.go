package service

import (
	"17mei/errs"
	"errors"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/store/db"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//AccessStoreRecylingInfo 获取云店回收信息
func (s *StoreServiceServer) AccessStoreRecylingInfo(ctx context.Context, in *pb.Recyling) (*pb.RecylingResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AccessStoreRecylingInfo", "%#v", in))

	err := db.AccessStoreRecylingInfo(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.RecylingResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UserSubmitRecylingOrder 提交预约订单接口
func (s *StoreServiceServer) UserSubmitRecylingOrder(ctx context.Context, in *pb.RecylingOrder) (*pb.RecylingOrderResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UserSubmitRecylingOrder", "%#v", in))

	err := db.UserSubmitRecylingOrder(in)

	if err != nil && err.Error() == "alreadyExists" {

		return &pb.RecylingOrderResp{Code: "00000", Message: "alreadyExists", Data: in}, nil

	} else if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.RecylingOrderResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UserAccessPendingRecylingOrder 查看预约中的回收订单接口
func (s *StoreServiceServer) UserAccessPendingRecylingOrder(ctx context.Context, in *pb.RecylingOrder) (*pb.RecylingOrderResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UserAccessPendingRecylingOrder", "%#v", in))

	err := db.UserAccessPendingRecylingOrder(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.RecylingOrderResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UpdateStoreRecylingInfo 设置云店回收信息
func (s *StoreServiceServer) UpdateStoreRecylingInfo(ctx context.Context, in *pb.Recyling) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateStoreRecylingInfo", "%#v", in))

	err := db.UpdateStoreRecylingInfo(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//GetStoreRecylingOrderList 获取云店回收订单列表
func (s *StoreServiceServer) GetStoreRecylingOrderList(ctx context.Context, in *pb.RecylingOrder) (*pb.RecylingOrderListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetStoreRecylingOrderList", "%#v", in))

	models, err, totalCount := db.GetStoreRecylingOrderList(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.RecylingOrderListResp{Code: "00000", Message: "ok", Data: models, TotalCount: totalCount}, nil
}

//UpdateRecylingOrder 更改回收订单
func (s *StoreServiceServer) UpdateRecylingOrder(ctx context.Context, in *pb.RecylingOrder) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateRecylingOrder", "%#v", in))

	err := db.UpdateRecylingOrder(in)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
