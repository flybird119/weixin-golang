package router

import m "github.com/goushuyun/weixin-golang/gateway/middleware"

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	return v1
}
