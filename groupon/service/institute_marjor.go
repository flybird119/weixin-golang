package service

import (
	"errors"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/groupon/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//通用专业批量增加
func (s *GrouponServiceServer) SharedMajorBatchSave(ctx context.Context, in *pb.SharedMajor) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SharedMajorBatchSave", "%#v", in))
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	defer tx.Rollback()
	for i := 0; i < len(in.Majors); i++ {
		if in.Majors[i] == nil {
			continue
		}
		err := db.SaveMarjor(tx, in.Majors[i])
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}
	tx.Commit()
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//获取专业列表（筛选获取）
func (s *GrouponServiceServer) SharedMajorList(ctx context.Context, in *pb.SharedMajor) (*pb.SharedMajorListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SharedMajorList", "%#v", in))
	majors, err, totalCount := db.FindMajorList(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SharedMajorListResp{Code: "00000", Message: "ok", TotalCount: totalCount, Data: majors}, nil
}

//创建学校的学院
func (s *GrouponServiceServer) SaveSchoolInstitute(ctx context.Context, in *pb.SchoolInstitute) (*pb.SchoolInstituteResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveSchoolInstitute", "%#v", in))
	err := db.SaveSchoolInstitute(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SchoolInstituteResp{Code: "00000", Message: "ok", Data: in}, nil
}

//创建学院专业
func (s *GrouponServiceServer) SaveInstituteMajor(ctx context.Context, in *pb.InstituteMajor) (*pb.InstituteMajorResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveInstituteMajor", "%#v", in))
	err := db.SaveInstituteMajor(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.InstituteMajorResp{Code: "00000", Message: "ok", Data: in}, nil
}

//获取学校学院专业列表
func (s *GrouponServiceServer) GetSchoolMajorInfo(ctx context.Context, in *pb.SchoolMajorInfoReq) (*pb.SchoolMajorListResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetSchoolMajorInfo", "%#v", in))
	schools, err := db.GetSchoolMajorInfo(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SchoolMajorListResp{Code: "00000", Message: "ok", Data: schools}, nil
}
