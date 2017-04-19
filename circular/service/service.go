package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/circular/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type CircularServiceServer struct{}

//AddCircular 增加轮播图
func (s *CircularServiceServer) AddCircular(ctx context.Context, in *pb.Circular) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddCircular", "%#v", in))
	err := db.AddCircular(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//DelCircular 删除轮播图
func (s *CircularServiceServer) DelCircular(ctx context.Context, in *pb.Circular) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "DelCircular", "%#v", in))
	err := db.DelCircular(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//UpdateCircular 修改轮播图信息
func (s *CircularServiceServer) UpdateCircular(ctx context.Context, in *pb.Circular) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateCircular", "%#v", in))
	err := db.UpdateCircular(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//CircularList 轮播图list
func (s *CircularServiceServer) CircularList(ctx context.Context, in *pb.Circular) (*pb.CircularListResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CircularList", "%#v", in))
	circulars, err := db.CircularList(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.CircularListResp{Code: "00000", Message: "ok", Data: circulars}, nil
}

//CircularList 初始化轮播图
func (s *CircularServiceServer) CircularInit(ctx context.Context, in *pb.Circular) (*pb.Void, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CircularInit", "%#v", in))
	//首先获取默认的所有的circular
	circulars, err := db.CircularList(&pb.Circular{StoreId: ""})
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	for i := 0; i < len(circulars); i++ {
		circular := circulars[i]

		circular.StoreId = in.StoreId
		circular.Profile = "默认简介"
		circular.Title = "默认标题"
		err := db.AddCircular(circular)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}
	return &pb.Void{}, nil
}
