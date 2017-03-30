package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1") // /v1/seller/test

	v1.Register("/seller/login", m.Wrap(c.SellerLogin))
	v1.Register("/seller/register", m.Wrap(c.SellerRegister))
	v1.Register("/seller/check_mobile", m.Wrap(c.CheckMobileExist))
	v1.Register("/seller/get_sms", m.Wrap(c.GetTelCode))

	return v1
}
