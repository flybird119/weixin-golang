package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc/token"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//AddStore 增加店铺接口
func AddStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "AddStore", req, "name")
}

//UpdateStore 增加店铺接口
func UpdateStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debugf("================%s", c.StoreId)

	storeid := c.StoreId
	log.Debugf("================%s", storeid)
	req := &pb.Store{Id: storeid, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "UpdateStore", req, "name", "profile")
}

//AddRealStore 增加实体店
func AddRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	req := &pb.RealStore{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "AddRealStore", req, "name", "province_code", "city_code", "scope_code", "address", "images")
}

//UpdateRealStore 修改实体店信息
func UpdateRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debugf("%+v", c)
	req := &pb.RealStore{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "UpdateRealStore", req, "name", "province_code", "city_code", "scope_code", "address", "images")
}

//StoreInfo 获取云店铺信息
func StoreInfo(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Id: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "StoreInfo", req)
}

//EnterStore 获取云店铺信息
func EnterStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "EnterStore", req, "id")
}

//ChangeStoreLogo 更改店铺logo
func ChangeStoreLogo(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Id: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "ChangeStoreLogo", req, "logo")
}

//RealStores 获取实体店列表
func RealStores(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Id: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "RealStores", req)
}
