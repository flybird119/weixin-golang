package service

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/wothing/log"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/school/db"
)

type SchoolServiceServer struct{}

//AddSchool 增加学校
func (s *SchoolServiceServer) AddSchool(ctx context.Context, in *pb.School) (*pb.SchoolResp, error) {
	//新建学校
	err := db.SaveSchool(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//记录日志
	log.Debugf("AddSchool SchoolId:%s Operater Id:%s", in.Id, in.Seller.Id)
	return &pb.SchoolResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UpdateSchool 更改学校基本信息
func (s *SchoolServiceServer) UpdateSchool(ctx context.Context, in *pb.School) (*pb.SchoolResp, error) {
	//更新学校
	err := db.UpdateSchool(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//记录日志
	log.Debugf("UpdateSchool Operation SchoolId:%s Operater Id:%s", in.Id, in.Seller.Id)
	return &pb.SchoolResp{Code: "00000", Message: "ok", Data: in}, nil
}

//UpdateExpressFee 更改运费
func (s *SchoolServiceServer) UpdateExpressFee(ctx context.Context, in *pb.School) (*pb.NormalResp, error) {
	//更改学校运费
	err := db.UpdateExpressFee(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	//记录日志
	log.Debugf("UpdateExpressFee Operation SchoolId:%s Operater Id:%s", in.Id, in.Seller.Id)
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//StoreSchools 店铺下的所有学校
func (s *SchoolServiceServer) StoreSchools(ctx context.Context, in *pb.School) (*pb.SchoolsResp, error) {
	//获取学校店铺
	schools, err := db.GetSchoolsByStore(in.StoreId)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SchoolsResp{Code: "00000", Message: "ok", Data: schools}, nil
}
