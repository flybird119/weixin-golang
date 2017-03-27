package router

import m "goushuyun/gateway/middleware"

func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	return v1
}
