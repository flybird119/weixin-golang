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
	log.Debug(c)
	req := &pb.Store{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "UpdateStore", req, "name", "profile", "id")
}

//AddRealStore 增加实体店
func AddRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	req := &pb.RealStore{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "AddRealStore", req, "name", "province_code", "city_code", "scope_code", "address", "images", "store_id")
}

//UpdateRealStore 修改实体店信息
func UpdateRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.RealStore{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "AddRealStore", req, "name", "province_code", "city_code", "scope_code", "address", "images", "id")
}

//StoreInfo 获取云店铺信息
func StoreInfo(w http.ResponseWriter, r *http.Request) {
	req := &pb.Store{}
	misc.CallWithResp(w, r, "bc_store", "StoreInfo", req, "id")
}

//StoreInfo 获取云店铺信息
func EnterStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "EnterStore", req, "id")
}
