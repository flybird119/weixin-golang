package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/misc/token"
	"github.com/goushuyun/weixin-golang/pb"
)

func AddGoods(w http.ResponseWriter, r *http.Request) {

	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}
	req := &pb.Goods{SellerId: c.SellerId, StoreId: c.StoreId}
	// call RPC to handle request
	misc.CallWithResp(w, r, "bc_goods", "AddGoods", req, "book_id", "isbn", "location")
}

func UpdateGoods(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}
	req := &pb.Goods{SellerId: c.SellerId, StoreId: c.StoreId}

	// call RPC to handle request
	misc.CallWithResp(w, r, "bc_goods", "UpdateGoods", req)
}

func SearchGoods(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}
	req := &pb.Goods{SellerId: c.SellerId, StoreId: c.StoreId}
	// call RPC to handle request
	misc.CallWithResp(w, r, "bc_goods", "SearchGoods", req)
}
