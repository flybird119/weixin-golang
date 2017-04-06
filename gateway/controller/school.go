package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc/token"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

//AddSchool 增加学校
func AddSchool(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		returnNotToken(w, r)
		return
	}
	req := &pb.School{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_school", "AddSchool", req, "name", "tel", "express_fee", "lat", "lng")
}

//UpdateSchool 更改学校基本信息
func UpdateSchool(w http.ResponseWriter, r *http.Request) {

	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		returnNotToken(w, r)
		return
	}
	req := &pb.School{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_school", "UpdateSchool", req, "id", "name", "tel", "express_fee", "lat", "lng")
}

//UpdateSchool 更改学校基本信息
func UpdateExpressFee(w http.ResponseWriter, r *http.Request) {
	//检测token
	c := token.Get(r)
	if c == nil || c.StoreId == "" {
		returnNotToken(w, r)
		return
	}
	req := &pb.School{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_school", "UpdateExpressFee", req, "id", "express_fee")
}

//StoreSchools 店铺下的所有学校
func StoreSchools(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		returnNotToken(w, r)
		return
	}
	req := &pb.School{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_school", "StoreSchools", req)
}

//returnNotToken 返回没找到token的错误提示
func returnNotToken(w http.ResponseWriter, r *http.Request) {

	misc.RespondMessage(w, r, map[string]interface{}{
		"code":    errs.ErrTokenNotFound,
		"message": "need token but not found",
	})
}
