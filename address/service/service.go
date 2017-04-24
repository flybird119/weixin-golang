package service

import (
	"17mei/errs"
	"errors"

	"github.com/goushuyun/weixin-golang/address/db"
	"github.com/goushuyun/weixin-golang/misc"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type AddressServiceServer struct{}

//增加地址
func (s *AddressServiceServer) AddAddress(ctx context.Context, in *pb.AddressReq) (*pb.AddressResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))

	//add operation
	for i := 0; i < len(in.Infos); i++ {
		in.Infos[i].UserId = in.UserId
		err := db.AddAddress(in.Infos[i])
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}

	return &pb.AddressResp{Code: "00000", Message: "ok", Data: in.Infos}, nil
}

//更新地址
func (s *AddressServiceServer) UpdateAddress(ctx context.Context, in *pb.AddressInfo) (*pb.AddressResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))
	//update operation
	err := db.UpdateAddress(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	infos := []*pb.AddressInfo{}
	infos = append(infos, in)
	return &pb.AddressResp{Code: "00000", Message: "ok", Data: infos}, nil
}

//我的地址
func (s *AddressServiceServer) MyAddresses(ctx context.Context, in *pb.AddressInfo) (*pb.AddressResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))
	infos, err := db.FindAddressByUser(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.AddressResp{Code: "00000", Message: "ok", Data: infos}, nil
}

//删除我的地址
func (s *AddressServiceServer) DeleteAddress(ctx context.Context, in *pb.AddressReq) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))
	err := db.DelAddress(in.Infos, in.UserId)

	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
