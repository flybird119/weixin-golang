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
	log.Debug("===================0")
	res, err, totalCount := db.SearchGoods(in)
	log.Debug("===================5")
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GoodsDetailResp{Code: "00000", Message: "ok", Data: res, TotalCount: totalCount}, nil
}

//GetGoodsByIdOrIsbn 获取商品基本信息
func (s *GoodsServiceServer) GetGoodsByIdOrIsbn(ctx context.Context, in *pb.Goods) (*pb.NormalGoodsResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetGoodsByIdOrIsbn", "%#v", in))
	err := db.GetGoodsByIdOrIsbn(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalGoodsResp{Code: "00000", Message: "ok", Data: in}, nil
}

//GetGoodsTypeInfo 获取商品单类型基础类型 主要用户购物车
func (s *GoodsServiceServer) GetGoodsTypeInfo(ctx context.Context, in *pb.TypeGoods) (*pb.TypeGoodsResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SearchGoodsTypeInfoByBookId", "%#v", in))
	//首先先获取商品的基本信息
	goods := &pb.Goods{Id: in.Id, StoreId: in.StoreId}
	err := db.GetGoodsByIdOrIsbn(goods)
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	//根据商品的信息 获取图书的信息
	data, err := misc.CallRPC(ctx, "bc_books", "GetBookInfo", &pb.Book{Id: goods.BookId})

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	book, ok := data.(*pb.Book)
	if !ok {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	in.BookId = goods.BookId
	in.StoreId = goods.StoreId
	in.Isbn = book.Isbn
	in.Title = book.Title
	in.Price = book.Price
	if in.Type == 0 {
		//新书
		in.Amount = goods.NewBookAmount
		in.SellingPrice = goods.NewBookPrice
		in.IsSelling = goods.HasNewBook

	} else {
		//二手书
		in.Amount = goods.OldBookAmount
		in.SellingPrice = goods.OldBookPrice
		in.IsSelling = goods.HasOldBook

	}
	in.Author = book.Author
	in.Publisher = book.Publisher
	in.GoodsImage = book.Image

	return &pb.TypeGoodsResp{Code: "00000", Message: "ok", Data: in}, nil
}

//DelOrOffShelfGoods 删除或者下架商品
func (s *GoodsServiceServer) DelOrRemoveGoods(ctx context.Context, in *pb.DelGoodsReq) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SearchGoodsTypeInfoByBookId", "%#v", in))

	//删除或者下架
	goodsModels := in.Data
	for i := 0; i < len(goodsModels); i++ {
		delModel := goodsModels[i]

		if delModel.GetOperateType() == 0 {
			//下架操作
			err := db.RemoveGoods(delModel)
			if err != nil {
				log.Debug(err)
				return nil, errs.Wrap(errors.New(err.Error()))
			}
		} else {
			//删除操作
			err := db.DelGoods(delModel)
			if err != nil {
				log.Debug(err)
				return nil, errs.Wrap(errors.New(err.Error()))
			}
		}

	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//GoodsLocationOperate 商品货架管理
func (s *GoodsServiceServer) GoodsLocationOperate(ctx context.Context, in *pb.GoodsLocation) (*pb.GoodsLocationResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GoodsLocationOperate", "%#v", in))

	//初始化数量
	in.Amount = 0
	if in.OperateType == 0 {
		//更新货架
		err := db.UpdateGoodsLocation(in)
		if err != nil {
			log.Debug(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}

	} else if in.OperateType == 1 {
		//删除货架
		err := db.DelGoodsLocation(in)
		if err != nil {
			log.Debug(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	} else {
		//增加货架
		err := db.InsertGoodsLocation(in)
		if err != nil {
			log.Debug(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}
	return &pb.GoodsLocationResp{Code: "00000", Message: "ok", Data: in}, nil
}

//AppSearchGoods 搜索图书 isbn 用于用户端搜索
func (s *GoodsServiceServer) AppSearchGoods(ctx context.Context, in *pb.Goods) (*pb.GoodsDetailResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GoodsLocationOperate", "%#v", in))

	res, err, totalCount := db.SearchGoodsNoLocation(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.GoodsDetailResp{Code: "00000", Message: "ok", Data: res, TotalCount: totalCount}, nil
}
