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
	if in.StoreId == "" {
		return nil, errs.Wrap(errors.New("需要重新登录！"))
	}
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
	schools, err := db.GetSchoolsByStore(in.StoreId, in.Status)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SchoolsResp{Code: "00000", Message: "ok", Data: schools}, nil
}

//StoreSchools 店铺下的所有学校
func (s *SchoolServiceServer) GetSchoolById(ctx context.Context, in *pb.School) (*pb.SchoolResp, error) {
	//获取学校店铺
	serchSchool, err := db.GetSchoolById(in.Id)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.SchoolResp{Code: "00000", Message: "ok", Data: serchSchool}, nil
}

//DelSchool 删除学校
func (s *SchoolServiceServer) DelSchool(ctx context.Context, in *pb.School) (*pb.NormalResp, error) {
	//获取学校店铺
	err := db.DelSchool(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//DelSchool 删除学校
func (s *SchoolServiceServer) UpdateSchoolRecylingState(ctx context.Context, in *pb.School) (*pb.NormalResp, error) {
	//获取学校店铺
	err := db.UpdateSchoolRecylingState(in)
	if err != nil {
		log.Debug(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
