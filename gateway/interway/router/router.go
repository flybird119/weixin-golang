package router

import (
	c "github.com/goushuyun/weixin-golang/gateway/controller"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

//SetRouterV1 设置seller的router
func SetRouterV1() *m.Router {
	v1 := m.NewWithPrefix("/v1") // /v1/seller/test
	//seller 开始
	v1.Register("/seller/login", m.Wrap(c.SellerLogin))
	v1.Register("/seller/register", m.Wrap(c.SellerRegister))
	v1.Register("/seller/check_mobile", m.Wrap(c.CheckMobileExist))
	v1.Register("/seller/get_sms", m.Wrap(c.GetTelCode))
	v1.Register("/seller/self_stores", m.Wrap(c.SelfStores))
	v1.Register("/seller/get_update_sms", m.Wrap(c.GetUpdateTelCode))
	v1.Register("/seller/update_password", m.Wrap(c.UpdatePasswordAndLogin))
	//store 开始
	v1.Register("/store/add", m.Wrap(c.AddStore))
	v1.Register("/store/update", m.Wrap(c.UpdateStore))
	v1.Register("/store/add_real_store", m.Wrap(c.AddRealStore))
	v1.Register("/store/update_real_store", m.Wrap(c.UpdateRealStore))
	v1.Register("/store/store_info", m.Wrap(c.StoreInfo))

	// mediastore
	v1.Register("/mediastore/get_upload_token", m.Wrap(c.GetUplaodToken))

	return v1
}
