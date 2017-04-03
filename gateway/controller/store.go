package controller

import (
	"17mei/misc/token"
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//AddStore 增加店铺接口
func AddStore(w http.ResponseWriter, r *http.Request) {
	req := &pb.Store{}
	misc.CallWithResp(w, r, "bc_store", "AddStore", req, "name")
}

//UpdateStore 增加店铺接口
func UpdateStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)
	req := &pb.Store{Seller: &pb.SellerInfo{}}
	misc.CallWithResp(w, r, "bc_store", "UpdateStore", req, "name", "profile", "id")
}
