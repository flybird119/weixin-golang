package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/retail/db"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type RetailServiceServer struct{}

//提交售后
func (s *RetailServiceServer) RetailSubmit(ctx context.Context, in *pb.RetailSubmitModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "RetailSubmit", "%#v", in))
	if len(in.Items) <= 0 {
		return nil, errs.Wrap(errors.New("items should not null"))
	}
	err := db.AddRetail(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//售后检索
func (s *RetailServiceServer) RetailList(ctx context.Context, in *pb.Retail) (*pb.RetailListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "RetailSubmit", "%#v", in))
	details, err, count := db.FindRetails(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.RetailListResp{Code: "00000", Message: "ok", Data: details, TotalCount: count}, nil
}
