package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc/token"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

//提交订单
func OrderSubmit(w http.ResponseWriter, r *http.Request) {
	req := &pb.OrderSubmitModel{}
	c := token.Get(r)
	// get store_id
	if c != nil && c.StoreId != "" && c.UserId != "" {
		req.StoreId = c.StoreId
		req.UserId = c.UserId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "OrderSubmit", req, "mobile", "name", "address", "school_id")
}

//模拟支付成功
func PaySuccess(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	misc.CallWithResp(w, r, "bc_order", "PaySuccess", req)
}

//app端订单列表
func OrderListApp(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	//搜索类型 来自用户
	c := token.Get(r)
	// get store_id
	if c != nil && c.StoreId != "" && c.UserId != "" {
		req.StoreId = c.StoreId
		req.UserId = c.UserId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "OrderList", req)
}

//seller端订单列表
func OrderListSeller(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	//搜索类型 来自商家
	req.SearchType = 1
	c := token.Get(r)

	// get store_id
	if c != nil && c.StoreId != "" && c.SellerId != "" {
		req.StoreId = c.StoreId
		req.SellerId = c.SellerId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "OrderList", req)
}

//打印订单
func PrintOrder(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	//搜索类型 来自商家
	c := token.Get(r)

	// get store_id
	if c != nil && c.StoreId != "" && c.SellerId != "" {
		req.StoreId = c.StoreId
		req.SellerId = c.SellerId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "PrintOrder", req)
}

//发货订单
func DeliverOrder(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	//搜索类型 来自商家
	c := token.Get(r)

	// get store_id
	if c != nil && c.StoreId != "" && c.SellerId != "" {
		req.StoreId = c.StoreId
		req.SellerId = c.SellerId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "DeliverOrder", req)
}

//配送订单
func DistributeOrder(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	//搜索类型 来自商家
	c := token.Get(r)

	// get store_id
	if c != nil && c.StoreId != "" && c.SellerId != "" {
		req.StoreId = c.StoreId
		req.SellerId = c.SellerId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "DistributeOrder", req)
}

//确认订单
func ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	req := &pb.Order{}
	//搜索类型 来自商家
	c := token.Get(r)

	// get store_id
	if c != nil && c.StoreId != "" && c.UserId != "" {
		req.StoreId = c.StoreId
		req.SellerId = c.SellerId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_order", "ConfirmOrder", req)
}
