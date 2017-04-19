package service

import (
	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type CircularServiceServer struct{}

//AddCircular 增加购物车
func (s *CircularServiceServer) AddCircular(ctx context.Context, in *pb.Circular) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddCircular", "%#v", in))

	// if err != nil {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//DelCircular 增加购物车
func (s *CircularServiceServer) DelCircular(ctx context.Context, in *pb.Circular) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "DelCircular", "%#v", in))

	// if err != nil {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//UpdateCircular 增加购物车
func (s *CircularServiceServer) UpdateCircular(ctx context.Context, in *pb.Circular) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateCircular", "%#v", in))

	// if err != nil {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//CircularList 增加购物车
func (s *CircularServiceServer) CircularList(ctx context.Context, in *pb.Circular) (*pb.CircularListResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CircularList", "%#v", in))

	// if err != nil {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }
	return &pb.CircularListResp{Code: "00000", Message: "ok"}, nil
}

//CircularList 增加购物车
func (s *CircularServiceServer) CircularInit(ctx context.Context, in *pb.Circular) (*pb.Void, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CircularList", "%#v", in))

	// if err != nil {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }
	return nil, nil
}
