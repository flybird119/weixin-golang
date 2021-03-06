package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1")

	// user
	v1.Register("/user/get_user_info", m.Wrap(c.GetUserInfo))

	// payment
	v1.Register("/payment/get_charge", m.Wrap(c.GetCharge))
	v1.Register("/payment/pay_success_notify", m.Wrap(c.PaySuccessNotify))
	v1.Register("/payment/refund_success_notify", m.Wrap(c.RefundSuccessNotify))

	// weixin
	v1.Register("/weixin/receive_verify_ticket", m.Wrap(c.ReceiveTicket))
	v1.Register("/weixin/get_weixin_info", m.Wrap(c.GetWeixinInfo))
	v1.Register("/weixin/get_js_ticket", m.Wrap(c.GetJsTicket))
	v1.Register("/weixin/extract_image_from_weixin_to_qiniu", m.Wrap(c.ExtractImg))
	v1.Register("/weixin/msgpush/:appid", m.Wrap(c.MsgPush))
	v1.Register("/weixin/get_user_base_info", m.Wrap(c.GetUserBaseInfo))
	v1.Register("/weixin/get_openid", m.Wrap(c.GetOpenid))
	v1.Register("/weixin/get_official_account_info", m.Wrap(c.GetOfficeAccountInfo))

	//school
	v1.Register("/school/get_store_schools", m.Wrap(c.StoreSchoolsApp))
	v1.Register("/school/get_school_info", m.Wrap(c.GetSchoolById))

	//goods
	v1.Register("/goods/search", m.Wrap(c.AppSearchGoods))

	//circular
	v1.Register("/circular/list", m.Wrap(c.CircularListApp))
	//cart
	v1.Register("/cart/add", m.Wrap(c.CartAdd))
	v1.Register("/cart/list", m.Wrap(c.CartList))
	v1.Register("/cart/update", m.Wrap(c.CartUpdate))
	v1.Register("/cart/del", m.Wrap(c.CartDel))
	//topic
	v1.Register("/topic/list", m.Wrap(c.TopicsInfoApp))
	//order
	v1.Register("/order/submit", m.Wrap(c.OrderSubmit))
	v1.Register("/order/pay_success", m.Wrap(c.PaySuccess))
	v1.Register("/order/find", m.Wrap(c.OrderListApp))
	v1.Register("/order/confirm", m.Wrap(c.ConfirmOrder))
	v1.Register("/order/detail", m.Wrap(c.OrderDetail))
	v1.Register("/order/after_sale_apply", m.Wrap(c.AfterSaleApply))
	v1.Register("/order/close", m.Wrap(c.CloseOrder))
	v1.Register("/order/necessary_order_counts", m.Wrap(c.UserCenterNecessaryOrderCount))
	v1.Register("/order/sharing_order_operation", m.Wrap(c.OrderShareOperation))

	//address
	v1.Register("/address/add", m.Wrap(c.AddAddress))
	v1.Register("/address/update", m.Wrap(c.UpdateAddress))
	v1.Register("/address/del", m.Wrap(c.DeleteAddress))
	v1.Register("/address/my_address", m.Wrap(c.MyAddresses))

	//store
	v1.Register("/store/info", m.Wrap(c.StoreInfoApp))
	v1.Register("/store/real_stores", m.Wrap(c.RealStores))

	//recyling
	v1.Register("/recyling/my_recyling_order", m.Wrap(c.UserAccessPendingRecylingOrder))
	v1.Register("/recyling/submit_recyling_order", m.Wrap(c.UserSubmitRecylingOrder))
	v1.Register("/recyling/store_recyling_info", m.Wrap(c.AccessStoreRecylingInfo))

	//groupon
	v1.Register("/groupon/get_school_majors", m.Wrap(c.GetSchoolMajorInfoApp))
	v1.Register("/groupon/save_groupon", m.Wrap(c.SaveGrouponApp))
	v1.Register("/groupon/find_groupon", m.Wrap(c.GrouponListApp))
	v1.Register("/groupon/my_groupon", m.Wrap(c.MyGrouponApp))
	v1.Register("/groupon/groupon_items", m.Wrap(c.GetGrouponItems))
	v1.Register("/groupon/groupon_related_user", m.Wrap(c.GetGrouponPurchaseUsers))
	v1.Register("/groupon/groupon_log", m.Wrap(c.GetGrouponOperateLog))
	v1.Register("/groupon/update_groupon", m.Wrap(c.UpdateGruopon))
	v1.Register("/groupon/star_groupon", m.Wrap(c.StarGroupon))
	v1.Register("/groupon/share_groupon", m.Wrap(c.ShareGroupon))
	v1.Register("/groupon/submit_groupon_order", m.Wrap(c.GrouponSubmit))
	v1.Register("/groupon/save_user_school_status", m.Wrap(c.SaveUserSchoolStatus))
	v1.Register("/groupon/update_user_school_status", m.Wrap(c.UpdateUserSchoolStatus))
	v1.Register("/groupon/get_user_school_status", m.Wrap(c.GetUserSchoolStatus))
	v1.Register("/groupon/user_has_star_groupon", m.Wrap(c.HasStarGroupon))
	return v1
}
