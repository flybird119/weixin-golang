package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"goushuyun/errs"

	"github.com/goushuyun/weixin-golang/misc"

	"github.com/goushuyun/weixin-golang/pb"
	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type PaymentService struct{}

func (s *PaymentService) GetCharge(ctx context.Context, req *pb.GetChargeReq) (*pb.GetChargeResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetCharge", "%#v", req))

	// 封装数据，并请求 charge 对象
	extra := make(map[string]interface{})
	extra["open_id"] = req.Openid

	params := &pingpp.ChargeParams{
		Order_no:  req.OrderNo,
		App:       pingpp.App{Id: "app_4qnjLOWXbDKSPmbb"},
		Amount:    uint64(req.Amount),
		Channel:   req.Channel,
		Currency:  "cny",
		Subject:   req.Subject,
		Body:      req.Body,
		Extra:     extra,
		Client_ip: req.Ip,
	}
	ch, err := charge.New(params)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	js, err := json.Marshal(ch)

	log.Debugf("The charge obj is %s", js)

	if err != nil {
		log.Terror(tid, err)
		return &pb.GetChargeResp{Code: errs.ErrInternal, Message: err.Error()}, nil
	}

	return &pb.GetChargeResp{Code: errs.Ok, Message: "OkHello", Data: fmt.Sprintf("%s", js)}, nil
}
