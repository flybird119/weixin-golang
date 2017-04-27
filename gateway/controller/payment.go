package controller

import (
	"encoding/json"
	"goushuyun/errs"
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

type payload struct {
	Data Data `json:"data"`
}
type Data struct {
	Object *pb.PaySuccessCallbackPayload `json:"object"`
}

func PaySuccessNotify(w http.ResponseWriter, r *http.Request) {
	log.Debugf("The response is : %s\n", r.Context().Value("body"))

	// callback string
	callback, ok := r.Context().Value("body").([]byte)
	if !ok {
		log.Error("interface to string error")
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "json unmarshal error",
		})
		return
	}

	// callback struct
	p := &pb.PaySuccessCallbackPayload{}
	data := Data{Object: p}
	obj := payload{
		Data: data,
	}

	// unmarshal
	err := json.Unmarshal(callback, &obj)
	if err != nil {
		log.Error(err)
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "json unmarshal error",
		})
		return
	}

	log.Debugf("The callback obj is %+v\n", obj)

	// 封装支付成功请求对象
	order := &pb.Order{
		Id:         p.OrderNo,
		TradeNo:    p.TransactionNo,
		PayChannel: p.Channel,
	}

	log.Debugf("The order obj is %+v\n", order)

	misc.CallWithResp(w, r, "bc_payment", "PaySuccessNotify", order, "id", "trade_no", "pay_channel")
}

func GetCharge(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	req := &pb.GetChargeReq{Ip: ip}
	misc.CallWithResp(w, r, "bc_payment", "GetCharge", req, "channel", "order_no", "amount", "subject", "body")
}
