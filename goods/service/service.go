package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

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
	if in.StoreId == "" {
		return nil, errs.Wrap(errors.New("没有关联店铺，请重试！"))
	}
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

//GetGoodsByIdOrIsbn 获取商品基本信息
func (s *GoodsServiceServer) GetGoodsByIdOrIsbn(ctx context.Context, in *pb.Goods) (*pb.NormalGoodsResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SearchGoods", "%#v", in))
	err := db.GetGoodsByIdOrIsbn(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalGoodsResp{Code: "00000", Message: "ok", Data: in}, nil
}

//GetGoodsTypeInfo 获取商品单类型基础类型
func (s *GoodsServiceServer) SearchGoodsTypeInfo(ctx context.Context, in *pb.TypeGoods) (*pb.TypeGoodsResp, error) {
	//首先先获取商品的基本信息

	return &pb.TypeGoodsResp{Code: "00000", Message: "ok"}, nil
}
