package service

import (
	"17mei/errs"
	"errors"

	"github.com/goushuyun/weixin-golang/goods/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//GoodsServiceServer service
type GoodsServiceServer struct{}

//AddGoods 增加商品
func (s *GoodsServiceServer) AddGoods(ctx context.Context, in *pb.Goods) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))
	err := db.AddGoods(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//UpdateGoods 更新商品
func (s *GoodsServiceServer) UpdateGoods(ctx context.Context, in *pb.Goods) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateGoods", "%#v", in))
	err := db.UpdateGoods(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//SearchGoods 查找商品
func (s *GoodsServiceServer) SearchGoods(ctx context.Context, in *pb.Goods) (*pb.GoodsDetailResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SearchGoods", "%#v", in))
	res, err := db.SearchGoods(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GoodsDetailResp{Code: "00000", Message: "ok", Data: res}, nil
}
