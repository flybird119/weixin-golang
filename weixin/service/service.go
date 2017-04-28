package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/misc/token"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/weixin/component"
	"github.com/goushuyun/weixin-golang/weixin/config"
	"github.com/goushuyun/weixin-golang/weixin/db"
	"github.com/wothing/log"
	"golang.org/x/net/context"

	"github.com/franela/goreq"
)

type WeixinServer struct{}

func (s *WeixinServer) WeChatJsApiTicket(ctx context.Context, req *pb.WeixinReq) (*pb.JsApiTicketResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "WeChatJsApiTicket", "%#v", req))

	// 获取对应公众号的 appid, refresh_token
	offical_account, err := db.GetAccountInfoByStoreId(req.StoreId)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 获取js_ticket
	ticket, err := component.JsTicket(offical_account.Appid, offical_account.RefreshToken)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	log.Debugf("<<<<<<<<<<<<<The Js_ticket is : %s>>>>>>>>>>>>>>", ticket)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := GetRandomString(16)

	text := fmt.Sprintf(`jsapi_ticket=%v&noncestr=%v&timestamp=%v&url=%v`, ticket, nonceStr, timestamp, req.Url)

	log.Debugf("需要加密的文本是：%s \n>>>>>>>>>>>>>>>>>>>>>", text)

	signature := Sha1Str(text)

	data := &pb.JsApiTicketResp_JsApiTicket{
		Appid:     offical_account.Appid,
		Signature: signature,
		Timestamp: timestamp,
		NonceStr:  nonceStr,
	}

	return &pb.JsApiTicketResp{Code: errs.Ok, Message: "ok", Data: data}, nil
}

func (s *WeixinServer) GetWeixinInfo(ctx context.Context, req *pb.WeixinReq) (*pb.GetWeixinInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetWeixinInfo", "%#v", req))

	// get weixin info
	weixinInfo, err := component.GetWeixinInfo(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// save weixin info to db
	userReq := &pb.User{WeixinInfo: weixinInfo, StoreId: req.StoreId}
	userResp := &pb.User{}
	err = misc.CallSVC(ctx, "bc_user", "SaveUser", userReq, userResp)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 修正 store_id 为当前云店铺的 store_id
	userResp.StoreId = req.StoreId

	// sign app token
	appToken := token.SignUserToken(token.AppToken, userResp.UserId, userResp.StoreId)

	return &pb.GetWeixinInfoResp{Code: errs.Ok, Message: "Ok", Data: userResp, Token: appToken}, nil
}

func (s *WeixinServer) GetOfficialAccountInfo(ctx context.Context, req *pb.WeixinReq) (*pb.NormalResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetOfficialAccountInfo", "%#v", req))

	// 使用auth_code去获取授权方的公众号帐号基本信息
	get_authorization_code_uri := "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s"

	type AuthorizationCodeItem struct {
		ComponentAppid    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}
	accessToken, err := component.ComponentAccessToken()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	config := config.GetConf()
	res, err := goreq.Request{
		Method: "POST",
		Uri:    fmt.Sprintf(get_authorization_code_uri, accessToken),
		Body: &AuthorizationCodeItem{
			ComponentAppid:    config.AppID,
			AuthorizationCode: req.AuthCode,
		},
	}.Do()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 解析res, 获取数据授权方 appid、authorizer_refresh_token , 并存入数据库
	GetApiQueryAuthResp := &pb.GetApiQueryAuth{}
	err = res.Body.FromJsonTo(GetApiQueryAuthResp)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	err = db.SaveAuthorizerInfoToStore(req.StoreId, GetApiQueryAuthResp.AuthorizationInfo.AuthorizerAppid, GetApiQueryAuthResp.AuthorizationInfo.AuthorizerRefreshToken)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 将API 授权 token 存入etcd
	err = saveAuthorizerAccessTokenToEtcd(GetApiQueryAuthResp.AuthorizationInfo.AuthorizerAppid, GetApiQueryAuthResp.AuthorizationInfo.AuthorizerRefreshToken)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 获取并存入微信公众号信息
	err = getandSaveAuthorizerAccountInfo(accessToken, config.AppID, GetApiQueryAuthResp.AuthorizationInfo.AuthorizerAppid)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: errs.Ok, Message: "ok"}, nil
}

func (s *WeixinServer) GetAuthURL(ctx context.Context, req *pb.WeixinReq) (*pb.GetAuthURLResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetAuthUrl", "%#v", req))

	conf := config.GetConf()
	pre_auth_code, err := component.PreAuthCode()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	// get auth url
	redirect_uri := "http://weixin.goushuyun.com/static/weixin_redirect.html"
	uri := "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"
	uri = fmt.Sprintf(uri, conf.AppID, pre_auth_code, redirect_uri)

	return &pb.GetAuthURLResp{Code: errs.Ok, Message: "ok", Url: uri}, nil
}
