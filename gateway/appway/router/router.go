package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	// weixin
	v1.Register("/weixin/receive_verify_ticket", m.Wrap(c.ReceiveTicket))
	v1.Register("/weixin/get_auth_url", m.Wrap(c.GetAuthURL))

	return v1
}
