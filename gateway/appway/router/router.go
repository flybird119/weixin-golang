package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	// payment
	v1.Register("/payment/get_charge", m.Wrap(c.GetCharge))

	// weixin
	v1.Register("/weixin/receive_verify_ticket", m.Wrap(c.ReceiveTicket))
	v1.Register("/weixin/get_weixin_info", m.Wrap(c.GetWeixinInfo))
	v1.Register("/weixin/test_db", m.Wrap(c.WeixinDBTest))
	//school
	v1.Register("/school/get_store_schools", m.Wrap(c.StoreSchools))
	v1.Register("/school/get_school_info", m.Wrap(c.GetSchoolById))

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
	v1.Register("/order/find", m.Wrap(c.OrderListApp))

	//address
	v1.Register("/address/add", m.Wrap(c.AddAddress))
	v1.Register("/address/update", m.Wrap(c.UpdateAddress))
	v1.Register("/address/del", m.Wrap(c.DeleteAddress))
	v1.Register("/address/my_address", m.Wrap(c.MyAddresses))

	return v1
}
