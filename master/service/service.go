package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/master/db"
	"github.com/goushuyun/weixin-golang/misc"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type MasterServiceServer struct{}

//管理员登陆
func (s *MasterServiceServer) MasterLogin(ctx context.Context, in *pb.Master) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "MasterLogin", "%#v", in))
	in.Id = ""
	err := db.GetUserByPasswordAndMobile(in)

	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if in.Id == "" {
		return &pb.NormalResp{Code: "00000", Message: "notFund"}, nil
	}
	//tokenStr := token.SignSellerToken(token.InterToken, sellerInfo.Id, sellerInfo.Mobile, "", 11)
	//sellerInfo.Token = tokenStr
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//提现列表
func (s *MasterServiceServer) WithdrawList(ctx context.Context, in *pb.StoreWithdrawalsModel) (*pb.WithdrawalsResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "WithdrawList", "%#v", in))
	models, err, totalCount := db.WithdrawList(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.WithdrawalsResp{Code: "00000", Message: "ok", TotalCount: totalCount, Data: models}, nil
}

//开始处理提现
func (s *MasterServiceServer) WithdrawHandle(ctx context.Context, in *pb.StoreWithdrawalsModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "WithdrawHandle", "%#v", in))
	model, err := db.GetWithdrawById(in.Id)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if model == nil {
		return nil, errs.Wrap(errors.New("bad parameter"))
	}
	if model.Status != 1 {
		return nil, errs.Wrap(errors.New("bad status"))
	}
	in.Status = 2
	in.AcceptAt = 1
	err = db.UpdateWithdraw(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}

//提现完成
func (s *MasterServiceServer) WithdrawComplete(ctx context.Context, in *pb.StoreWithdrawalsModel) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "WithdrawComplete", "%#v", in))
	model, err := db.GetWithdrawById(in.Id)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	if model == nil {
		return nil, errs.Wrap(errors.New("bad parameter"))
	}
	if model.Status != 1 {
		return nil, errs.Wrap(errors.New("bad status"))
	}
	in.Status = 3
	in.CompleteAt = 1
	err = db.UpdateWithdraw(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.NormalResp{Code: "00000", Message: "ok"}, nil
}
