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
	"github.com/pingplusplus/pingpp-go/pingpp/refund"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type PaymentService struct{}

func (s *PaymentService) Refund(ctx context.Context, req *pb.RefundReq) (*pb.Void, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "Refund", "%#v", req))

	params := &pingpp.RefundParams{
		Amount:      uint64(req.Amount),
		Description: req.Reason,
	}
	re, err := refund.New(req.TradeNo, params)
	if err != nil {
		// 生成退款请求时出错，可能为该笔订单已经足额退款
		log.Error(err)

		// 封装错误信息并返回
		callback := &pb.RefundErrCallback{}
		err = json.Unmarshal([]byte(err.Error()), callback)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		log.Debugf("********订单已足额退款*******\n %+v", callback)

		return nil, errors.New(callback.Message)
	}
	if !re.Succeed {
		// 退款失败，可能原因为商户平台余额不足

		log.Debugf("*********商户平台余额不足***********\n%+v", re)

		return nil, errors.New(re.Failure_msg)
	}

	log.Debug("*******退款完成，没毛病*********")

	return &pb.Void{}, nil
}

func (s *PaymentService) PaySuccessNotify(ctx context.Context, req *pb.Order) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetCharge", "%#v", req))

	// call rpc to notify pay success
	_, err := misc.CallRPC(ctx, "bc_order", "PaySuccess", req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: errs.Ok, Message: "ok"}, nil
}

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
