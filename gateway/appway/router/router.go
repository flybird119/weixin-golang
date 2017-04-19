package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	// user

	// weixin
	v1.Register("/weixin/receive_verify_ticket", m.Wrap(c.ReceiveTicket))
	v1.Register("/weixin/get_weixin_info", m.Wrap(c.GetWeixinInfo))

	//school
	v1.Register("/school/get_store_schools", m.Wrap(c.StoreSchools))

	//goods
	v1.Register("/goods/search", m.Wrap(c.AppSearchGoods))

	//circular
	v1.Register("/circular/list", m.Wrap(c.CircularList))
	//cart

	//topic
	v1.Register("/topic/list", m.Wrap(c.TopicsInfo))
	return v1
}
