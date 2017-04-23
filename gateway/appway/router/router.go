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
	v1.Register("/cart/add", m.Wrap(c.CartAdd))
	v1.Register("/cart/list", m.Wrap(c.CartList))
	v1.Register("/cart/update", m.Wrap(c.CartUpdate))
	v1.Register("/cart/del", m.Wrap(c.CartDel))
	//topic
	v1.Register("/topic/list", m.Wrap(c.TopicsInfo))
	//order
	v1.Register("/order/submit", m.Wrap(c.OrderSubmit))
	v1.Register("/order/pay_success", m.Wrap(c.PaySuccess))
	return v1
}
