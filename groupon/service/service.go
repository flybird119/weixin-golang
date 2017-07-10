package service

import (
	"17mei/errs"
	"errors"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/groupon/db"
	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//GrouponService service
type GrouponServiceServer struct{}

//创建班级购
func (s *GrouponServiceServer) SaveGroupon(ctx context.Context, in *pb.Groupon) (*pb.GrouponResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveGroupon", "%#v", in))
	err := db.SaveGroupon(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GrouponResp{Code: "00000", Message: "ok", Data: in}, nil
}

//班级购列表
func (s *GrouponServiceServer) GrouponList(ctx context.Context, in *pb.Groupon) (*pb.GrouponListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GrouponList", "%#v", in))
	models, err, totalCount := db.FindGroupon(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GrouponListResp{Code: "00000", Message: "ok", Data: models, TotalCount: totalCount}, nil
}

//班级购列表
func (s *GrouponServiceServer) MyGroupon(ctx context.Context, in *pb.Groupon) (*pb.GrouponListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "MyGroupon", "%#v", in))
	var models []*pb.Groupon
	var err error
	var totalCount int64
	if in.SearchOperateType == 1 {
		in.ParticipateUser = in.FounderId
		in.FounderId = ""
		models, err, totalCount = db.FindGroupon(in)
	} else {
		models, err, totalCount = db.FindGroupon(in)
	}
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GrouponListResp{Code: "00000", Message: "ok", Data: models, TotalCount: totalCount}, nil
}

//新增班级购操作日志
func (s *GrouponServiceServer) GetGrouponItems(ctx context.Context, in *pb.Groupon) (*pb.GrouponItemListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetGrouponItems", "%#v", in))
	models, err := db.GetGrouponItems(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GrouponItemListResp{Code: "00000", Message: "ok", Data: models}, nil

}

//获取班级购参与人信息
func (s *GrouponServiceServer) GetGrouponPurchaseUsers(ctx context.Context, in *pb.Groupon) (*pb.GrouponUserListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetGrouponPurchaseUsers", "%#v", in))
	models, err := db.GetGrouponPurchaseUsers(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GrouponUserListResp{Code: "00000", Message: "ok", Data: models}, nil
}

//获取班级购操作日志
func (s *GrouponServiceServer) GetGrouponOperateLog(ctx context.Context, in *pb.Groupon) (*pb.GrouponOperateLogListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetGrouponOperateLog", "%#v", in))
	models, err := db.GetGrouponOperateLog(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.GrouponOperateLogListResp{Code: "00000", Message: "ok", Data: models}, nil
}

//修改班级购
func (s *GrouponServiceServer) UpdateGruopon(ctx context.Context, in *pb.Groupon) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateGruopon", "%#v", in))
	err := db.UpdateGruopon(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	oplog := &pb.GrouponOperateLog{GrouponId: in.Id, FounderId: in.FounderId, FounderName: in.FounderName, FounderType: in.FounderType, OperateType: "update", OperateDetail: " "}
	err = db.SaveGrouponOperateLog(oplog)
	if err != nil {
		log.Error(err)
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//批量班级购日期 ---------------------------------------------------------------------
func (s *GrouponServiceServer) BatchUpdateGrouponExpireAt(ctx context.Context, in *pb.Groupon) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "BatchUpdateGrouponExpireAt", "%#v", in))
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//批量班级购日期
func (s *GrouponServiceServer) StarGroupon(ctx context.Context, in *pb.GrouponOperateLog) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "StarGroupon", "%#v", in))
	totalCount, _ := db.HasGrouponLogWithOpreation(in.GrouponId, in.FounderId, "star")
	if totalCount <= 0 {
		err := db.SaveGrouponOperateLog(in)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//点赞
func (s *GrouponServiceServer) ShareGroupon(ctx context.Context, in *pb.GrouponOperateLog) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ShareGroupon", "%#v", in))
	err := db.SaveGrouponOperateLog(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//下单
func (s *GrouponServiceServer) GrouponSubmit(ctx context.Context, in *pb.GrouponSubmitModel) (*pb.OrderSubmitResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GrouponSubmit", "%#v", in))

	order, noStack, err := db.GrouponSubmit(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//库存不足
	if noStack != "" {
		return &pb.OrderSubmitResp{Code: "00000", Message: "noStack", Data: order}, nil
	}
	return &pb.OrderSubmitResp{Code: "00000", Message: "ok", Data: order}, nil
}

//保存学生学籍信息
func (s *GrouponServiceServer) SaveUserSchoolStatus(ctx context.Context, in *pb.UserSchoolStatus) (*pb.UserSchoolStatusResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveUserSchoolStatus", "%#v", in))
	err := db.SaveUserSchoolStatus(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if in.Id == "" {
		return &pb.UserSchoolStatusResp{Code: "00000", Message: "exists", Data: in}, nil
	}
	return &pb.UserSchoolStatusResp{Code: "00000", Message: "ok", Data: in}, nil
}

//更新学生学籍信息
func (s *GrouponServiceServer) UpdateUserSchoolStatus(ctx context.Context, in *pb.UserSchoolStatus) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "UpdateUserSchoolStatus", "%#v", in))
	err := db.UpdateUserSchoolStatus(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//获取学生学籍
func (s *GrouponServiceServer) GetUserSchoolStatus(ctx context.Context, in *pb.UserSchoolStatus) (*pb.UserSchoolStatusResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetUserSchoolStatus", "%#v", in))
	err := db.GetUserSchoolStatus(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.UserSchoolStatusResp{Code: "00000", Message: "ok"}, nil
}
