package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/store/db"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

//获取提现账号管理的手机验证码
func (s *StoreServiceServer) RechargeApply(ctx context.Context, in *pb.RechargeModel) (*pb.RechargeOperationResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetCardOperSmsCode", "%#v", in))
	err := db.RechargeApply(in)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.RechargeOperationResp{Code: "00000", Message: "ok", Data: in}, nil
}
