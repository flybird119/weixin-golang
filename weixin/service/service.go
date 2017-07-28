package service

import (
	"database/sql"
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

	global_db "github.com/goushuyun/weixin-golang/db"

	"github.com/franela/goreq"
)

type WeixinServer struct{}

func (s WeixinServer) GetUserBaseInfo(ctx context.Context, req *pb.WeixinReq) (*pb.GetUserBaseInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetUserBaseInfo", "%#v", req))

	// 获取store_id 所对应的 appid
	// officialAccount, err := db.GetAccountInfoByStoreId(req.StoreId)
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, errs.Wrap(errors.New(err.Error()))
	// }

	// 拿到该 appid 的授权 access_token
	// authorizer_token, err := component.ApiAuthorizerToken(officialAccount.Appid, officialAccount.RefreshToken)

	// while get user's information, use [qudianchi] for all
	refresh_token := global_db.GetValue("weixin", "qudianchi_refresh_token", "refreshtoken@@@gpS28RSEBu16UiFm6yMzRt4-mQ9Pz_aBsdY4psYHOh4")
	authorizer_token, err := component.ApiAuthorizerToken("wx6d36779ce4dd3dfa", refresh_token)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	weixin_info, err := getUserBaseInfo(req.Openid, authorizer_token)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	log.Debugf("【get user info】 , openid : %s\n", req.Openid)
	log.JSONIndent(weixin_info)

	// 更新weixin_info
	if weixin_info.Subscribe == 0 {
		// 未订阅
		return &pb.GetUserBaseInfoResp{Code: errs.Ok, Message: "no_subscribe"}, nil
	} else {
		err = db.UpdateUserInfo(weixin_info, req.UserId)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}

	user := &pb.User{
		WeixinInfo: weixin_info,
	}

	// 请求 用户数据
	go func() {
		err = db.CreateUser2StoreMap(req)
		if err != nil {
			log.Error(err)
			// return nil, errs.Wrap(errors.New(err.Error()))
		}
	}()
	return &pb.GetUserBaseInfoResp{Code: errs.Ok, Message: "ok", Data: user}, nil
}

func (s *WeixinServer) GetOpenid(ctx context.Context, req *pb.GetUserInfoReq) (*pb.WeixinInfo, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetOpenid", "%#v", req))

	official_openid, err := component.GetOpenid(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.WeixinInfo{Openid: official_openid}, nil
}

func (s *WeixinServer) GetOfficeAccountInfo(ctx context.Context, req *pb.WeixinReq) (*pb.GetOfficeAccountInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetOfficeAccountInfo", "%#v", req))

	// 通过store_id, 取出 official_account 信息
	oa, err := db.GetAccountInfoByStoreId(req.StoreId)

	if err == sql.ErrNoRows {
		// 没有该 store_id 对应的 office_account
		log.Error(err)
		return &pb.GetOfficeAccountInfoResp{Code: errs.Ok, Message: "not_found"}, nil
	}

	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	oa.StoreId = req.StoreId
	return &pb.GetOfficeAccountInfoResp{Code: errs.Ok, Message: "ok", Data: oa}, nil
}

func (s *WeixinServer) ExtractImageFromWeixin(ctx context.Context, req *pb.ExtractImageReq) (*pb.ExtractImageResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "ExtraImageFromWeixin", "%#v", req))

	// 根据store_id 获取店铺的相关信息
	official_account, err := db.GetAccountInfoByStoreId(req.StoreId)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 获取对应appid的 Authorization access_token
	authorizer_token, err := component.ApiAuthorizerToken(official_account.Appid, official_account.RefreshToken)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// 封装 weixin urls
	weixin_media_urls := []string{}
	for _, server_id := range req.ServerIds {
		url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s", authorizer_token, server_id)
		weixin_media_urls = append(weixin_media_urls, url)
	}

	// 调用七牛抓取图片，并返回keys
	extract_req := &pb.ExtractReq{Appid: official_account.Appid, Zone: req.Zone, WeixinMediaUrls: weixin_media_urls}
	extract_resp := &pb.ExtractResp{}
	err = misc.CallSVC(ctx, "bc_mediastore", "ExtractImageFromWeixinToQiniu", extract_req, extract_resp)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.ExtractImageResp{Code: errs.Ok, Message: "ok", QiniuKeys: extract_resp.QiniuKeys}, nil
}

func (s *WeixinServer) WeChatJsApiTicket(ctx context.Context, req *pb.WeixinReq) (*pb.JsApiTicketResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "WeChatJsApiTicket", "%#v", req))

	// 获取对应公众号的 appid, refresh_token
	offical_account, err := db.GetAccountInfoByStoreId(req.StoreId)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	log.Debugf("The offical_account info is : %s", offical_account)

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

	log.Debug("------------------------------------")
	log.JSONIndent(weixinInfo)
	log.Debug("------------------------------------")

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
	err = saveAuthorizerAccessTokenToEtcd(GetApiQueryAuthResp.AuthorizationInfo.AuthorizerAppid, GetApiQueryAuthResp.AuthorizationInfo.AuthorizerAccessToken)
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
	uri := "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"
	uri = fmt.Sprintf(uri, conf.AppID, pre_auth_code, req.RedirectUri)

	return &pb.GetAuthURLResp{Code: errs.Ok, Message: "ok", Url: uri}, nil
}
