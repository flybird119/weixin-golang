package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/topic/db"
	"github.com/wothing/log"
)

type TopicServiceServer struct{}

//AddTopic 增加话题
func (s *TopicServiceServer) AddTopic(ctx context.Context, in *pb.Topic) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddTopic", "%#v", in))
	//****
	//    限制最多创建20个话题
	//*****
	seachTopic := &pb.Topic{TokenStoreId: in.TokenStoreId}
	topics, err, _ := db.SearchTopics(seachTopic)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if len(topics) > 20 {

		return nil, errs.Wrap(errors.New("已达创建上限"))
	}

	err = db.AddTopic(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//DelTopic 删除话题
func (s *TopicServiceServer) DelTopic(ctx context.Context, in *pb.Topic) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "DelTopic", "%#v", in))
	err := db.DelTopic(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//UpdateTopic 更新话题
func (s *TopicServiceServer) UpdateTopic(ctx context.Context, in *pb.Topic) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateTopic", "%#v", in))
	err := db.UpdateTopic(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//AddTopicItem 增加话题项
func (s *TopicServiceServer) AddTopicItem(ctx context.Context, in *pb.TopicItem) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddTopicItem", "%#v", in))
	err := db.AddTopicItem(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//DelTopicItem 删除话题项
func (s *TopicServiceServer) DelTopicItem(ctx context.Context, in *pb.TopicItem) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "DelTopicItem", "%#v", in))
	err := db.DelTopicItem(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//SearchTopics 搜索话题
func (s *TopicServiceServer) SearchTopics(ctx context.Context, in *pb.Topic) (*pb.SearchTopicResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SearchTopics", "%#v", in))
	topics, err, totalCount := db.SearchTopics(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SearchTopicResp{Code: "00000", Message: "ok", Data: topics, TotalCount: totalCount}, nil
}

//SearchTopics 搜索话题
func (s *TopicServiceServer) TopicsInfo(ctx context.Context, in *pb.Topic) (*pb.SearchTopicResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SearchTopics", "%#v", in))
	topics, err, totalCount := db.SearchTopics(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	for i := 0; i < len(topics); i++ {
		topic := topics[i]
		for j := 0; j < len(topic.Items); j++ {
			item := topic.Items[j]
			goods := &pb.Goods{StoreId: in.TokenStoreId, Id: item.GoodsId}
			log.Debugf("++++++++++++++++++++++++++++++%+v", goods)
			data, err := misc.CallRPC(ctx, "bc_goods", "GetGoodsByIdOrIsbn", goods)
			log.Debugf("++++++++++++++++++++++++++++++%#v", data)
			if err != nil {
				log.Debug(err)
				return nil, errs.Wrap(errors.New(err.Error()))
			}
			goodsResp, ok := data.(*pb.NormalGoodsResp)
			if !ok {
				log.Debug(err)
				return nil, errs.Wrap(errors.New(err.Error()))
			}
			data, err = misc.CallRPC(ctx, "bc_books", "GetBookInfo", &pb.Book{Id: goodsResp.Data.BookId})
			if err != nil {
				log.Debug(err)
				return nil, errs.Wrap(errors.New(err.Error()))
			}
			book, ok := data.(*pb.Book)
			if !ok {
				log.Debug(err)
				return nil, errs.Wrap(errors.New(err.Error()))
			}
			item.Isbn = book.Isbn
			item.Title = book.Title
			item.Author = book.Author
			item.Publisher = book.Publisher
			item.Image = book.Image
			item.BookPrice = book.Price
			item.Stock = goodsResp.Data.NewBookAmount + goodsResp.Data.OldBookAmount

		}
	}
	return &pb.SearchTopicResp{Code: "00000", Message: "ok", Data: topics, TotalCount: totalCount}, nil
}
