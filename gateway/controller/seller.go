package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

//SellerLogin 手机号登录
func SellerLogin(w http.ResponseWriter, r *http.Request) {
	req := &pb.LoginModel{}
	misc.CallWithResp(w, r, "bc_seller", "SellerLogin", req, "mobile", "password")
}

//SellerRegister 商家注册
func SellerRegister(w http.ResponseWriter, r *http.Request) {
	req := &pb.RegisterModel{}
	misc.CallWithResp(w, r, "bc_seller", "SellerRegister", req, "mobile", "password", "message_code", "username")
}

//CheckMobileExist 检验手机号是否注册过
func CheckMobileExist(w http.ResponseWriter, r *http.Request) {
	req := &pb.CheckMobileReq{}
	misc.CallWithResp(w, r, "bc_seller", "CheckMobileExist", req, "mobile")
}

//GetTelCode 获取手机验证码
func GetTelCode(w http.ResponseWriter, r *http.Request) {
	req := &pb.CheckMobileReq{}
	misc.CallWithResp(w, r, "bc_seller", "GetTelCode", req, "mobile")
}
