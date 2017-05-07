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

//UpdateStore 修改
func UpdateStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	storeid := c.StoreId
	req := &pb.Store{Id: storeid, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "UpdateStore", req, "name", "profile")
}

//AddRealStore 增加实体店
func AddRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.RealStore{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "AddRealStore", req, "name", "province_code", "city_code", "scope_code", "address", "images")
}

//UpdateRealStore 修改实体店信息
func UpdateRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.RealStore{Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "UpdateRealStore", req, "name", "province_code", "city_code", "scope_code", "address", "images")
}

//StoreInfo 获取云店铺信息
func StoreInfo(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

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

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.Store{Id: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "ChangeStoreLogo", req, "logo")
}

//RealStores 获取实体店列表
func RealStores(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.Store{Id: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "RealStores", req)
}

//检查code
func CheckCode(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.RegisterModel{}
	misc.CallWithResp(w, r, "bc_store", "CheckCode", req, "mobile", "message_code")
}

//TransferStore 转让店铺
func TransferStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.TransferStoreReq{Store: &pb.Store{Id: c.StoreId}}
	misc.CallWithResp(w, r, "bc_store", "TransferStore", req, "mobile", "message_code")
}

//删除实体店
func DelRealStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	log.Debug(c)

	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.RealStore{StoreId: c.StoreId, Seller: &pb.SellerInfo{Id: c.SellerId, Mobile: c.Mobile}}
	misc.CallWithResp(w, r, "bc_store", "DelRealStore", req, "id")

}

//删除实体店
func GetCardOperSmsCode(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.SmsCardSubmitModel{StoreId: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "GetCardOperSmsCode", req)

}

//保存提现账号
func SaveStoreWithdrawCard(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.StoreWithdrawCard{StoreId: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "SaveStoreWithdrawCard", req, "card_no", "card_name", "username", "code")

}

//保存提现账号
func UpdateStoreWithdrawCard(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}

	req := &pb.StoreWithdrawCard{StoreId: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "UpdateStoreWithdrawCard", req, "id")

}

//保存提现账号
func GetWithdrawCardInfoByStore(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}
	req := &pb.StoreWithdrawCard{StoreId: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "GetWithdrawCardInfoByStore", req)

}

//店铺首页历史订单各个状态统计
func StoreHistoryStateOrderNum(w http.ResponseWriter, r *http.Request) {
	c := token.Get(r)
	//检测token
	if c == nil || c.StoreId == "" {
		misc.ReturnNotToken(w, r)
		return
	}
	req := &pb.StoreHistoryStateOrderNumModel{StoreId: c.StoreId}
	misc.CallWithResp(w, r, "bc_store", "StoreHistoryStateOrderNum", req)

}
