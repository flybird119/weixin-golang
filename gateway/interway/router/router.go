package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	v1.Register("/seller/test", m.Wrap(c.SellerLogin))

	return v1
}
