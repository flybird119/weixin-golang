package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")
	//登陆
	v1.Register("/master/login", m.Wrap(c.MasterLogin))
	v1.Register("/master/withdraw_list", m.Wrap(c.WithdrawList))
	v1.Register("/master/withdraw_handle", m.Wrap(c.WithdrawHandle))
	v1.Register("/master/withdraw_complete", m.Wrap(c.WithdrawComplete))

	return v1
}
