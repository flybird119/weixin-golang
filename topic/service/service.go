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
	topics, err := db.SearchTopics(seachTopic)

	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if len(topics) > 20 {

		return nil, errs.Wrap(errors.New("已达创建上限"))
	}

	//******
	//  每个话题最多限制15本书
	//******
	if len(in.Items) > 15 {
		return nil, errs.Wrap(errors.New("已超话题最大商品数量"))
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
	topics, err := db.SearchTopics(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SearchTopicResp{Code: "00000", Message: "ok", Data: topics}, nil
}
