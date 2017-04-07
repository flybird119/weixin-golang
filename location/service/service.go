package service

import (
	"errors"
	"goushuyun/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/location/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type LocationServiceServer struct{}

func (s *LocationServiceServer) GetChildrenLocation(ctx context.Context, req *pb.Location) (*pb.GetChildrenLocationResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetChildrenLocation", "%#v", req))

	// to get location's children
	err := db.GetDescLocation(req, req.Level)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.GetChildrenLocationResp{Code: errs.Ok, Message: "ok", Data: req.Children}, nil
}

func (s *LocationServiceServer) ListLocation(ctx context.Context, req *pb.Location) (*pb.ListLocationResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ListLocation", "%#v", req))

	// list location
	data, err := db.ListLocation(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.ListLocationResp{Code: errs.Ok, Message: "ok", Data: data}, nil
}

func (s *LocationServiceServer) UpdateLocation(ctx context.Context, req *pb.Location) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateLocation", "%#v", req))

	// update location, including name、pid
	if err := db.UpdateLocation(req); err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: errs.Ok, Message: "ok"}, nil
}

func (s *LocationServiceServer) AddLocation(ctx context.Context, req *pb.Location) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddLocation", "%#v", req))

	// insert location into db
	if err := db.AddLocation(req); err != nil {
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: errs.Ok, Message: "ok"}, nil
}
