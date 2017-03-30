package controller

import (
	"fmt"
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

func SellerLogin(w http.ResponseWriter, r *http.Request) {
	req := &pb.LoginModel

	fmt.Fprintf(w, "Hello, world")

	misc.CallWithResp(w, r, "seller", "SellerLogin", req, "mobile", "password")
}
