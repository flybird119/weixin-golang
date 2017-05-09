package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	// payment
	v1.Register("/master/get_charge", m.Wrap(c.GetCharge))
	v1.Register("/master/pay_success_notify", m.Wrap(c.PaySuccessNotify))
	v1.Register("/master/refund_success_notify", m.Wrap(c.RefundSuccessNotify))
	return v1
}
