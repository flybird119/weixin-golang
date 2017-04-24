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
		err := db.AddAddress(in.Infos[i])
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}

	return &pb.AddressResp{}, nil
}

//更新地址
func (s *AddressServiceServer) UpdateAddress(ctx context.Context, in *pb.AddressInfo) (*pb.AddressResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))

	return &pb.AddressResp{}, nil
}

//我的地址
func (s *AddressServiceServer) MyAddresses(ctx context.Context, in *pb.AddressInfo) (*pb.AddressResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))

	return &pb.AddressResp{}, nil
}

//删除我的地址
func (s *AddressServiceServer) DeleteAddress(ctx context.Context, in *pb.AddressReq) (*pb.NormalResp, error) {

	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "AddGoods", "%#v", in))

	return &pb.NormalResp{}, nil
}
