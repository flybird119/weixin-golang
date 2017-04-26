package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

func GetCharge(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	req := &pb.GetChargeReq{Ip: ip}
	misc.CallWithResp(w, r, "bc_payment", "GetCharge", req, "channel", "order_no", "amount", "subject", "body")
}
