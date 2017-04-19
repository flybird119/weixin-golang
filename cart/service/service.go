package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/cart/db"
	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//CartServiceServer 结构体
type CartServiceServer struct{}

//CartAdd 增加购物车
func (s *CartServiceServer) CartAdd(ctx context.Context, in *pb.Cart) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CartAdd", "%#v", in))
	err := db.CartAdd(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//CartList 购物车列表
func (s *CartServiceServer) CartList(ctx context.Context, in *pb.Cart) (*pb.CartListResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CartList", "%#v", in))
	carts, err := db.CartList(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//CallRPC
	for i := 0; i < len(carts); i++ {
		req := &pb.TypeGoods{Id: in.GoodsId, Type: in.Type}
		data, err := misc.CallRPC(ctx, "bc_goods", "GetGoodsTypeInfo", req)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}

		typeGoodsResp, ok := data.(*pb.TypeGoodsResp)
		if !ok {
			log.Error("断言失败")
			return nil, errs.Wrap(errors.New("something is error"))
		}
		typeGoods := typeGoodsResp.Data
		carts[i].GoodsDetail = typeGoods

	}
	return &pb.CartListResp{Code: "00000", Message: "ok", Data: carts}, nil
}

//CartUpdate 更改购物车
func (s *CartServiceServer) CartUpdate(ctx context.Context, in *pb.CartUpdateReq) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CartUpdate", "%#v", in))
	carts := in.Carts
	for i := 0; i < len(carts); i++ {
		err := db.CartUpdate(carts[i])
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//CartDel 删除购物车
func (s *CartServiceServer) CartDel(ctx context.Context, in *pb.Cart) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "CartDel", "%#v", in))
	err := db.CartDel(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
