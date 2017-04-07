package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/misc/token"
	"github.com/goushuyun/weixin-golang/pb"
)

func GetChildrenLocation(w http.ResponseWriter, r *http.Request) {
	req := &pb.Location{}
	misc.CallWithResp(w, r, "bc_location", "GetChildrenLocation", req, "id", "name")
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	req := &pb.Location{}
	misc.CallWithResp(w, r, "bc_location", "UpdateLocation", req, "id")
}

func ListLocation(w http.ResponseWriter, r *http.Request) {
	req := &pb.Location{}
	if c := token.Get(r); c != nil {
		req.StoreId = c.StoreId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}

	// call RPC to handle request
	misc.CallWithResp(w, r, "bc_location", "ListLocation", req)
}

func AddLocation(w http.ResponseWriter, r *http.Request) {
	req := &pb.Location{}
	if c := token.Get(r); c != nil {
		req.StoreId = c.StoreId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}

	// call RPC to handle request
	misc.CallWithResp(w, r, "bc_location", "AddLocation", req, "name")
}
