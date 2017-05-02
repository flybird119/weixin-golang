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
func (s *RetailServiceServer) RetailSubmit(ctx context.Context, in *pb.RetailSubmitModel) (*pb.RetailSubmitResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "RetailSubmit", "%#v", in))
	if len(in.Items) <= 0 {
		return nil, errs.Wrap(errors.New("items should not null"))
	}
	err := db.AddRetail(in)
	log.Debug("=============")
	log.Debug(err.Error())
	log.Debug("=============")
	if err.Error() == "noStock" {
		return &pb.RetailSubmitResp{Code: "00000", Message: "noStock", Items: in.Items}, nil
	}
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.RetailSubmitResp{Code: "00000", Message: "ok", Items: in.Items}, nil
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
