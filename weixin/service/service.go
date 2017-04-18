package service

import (
	"errors"
	"fmt"

	"github.com/goushuyun/weixin-golang/misc"

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

	// 解析res, 获取数据授权方 appid, 并存入数据库
	GetApiQueryAuthResp := &pb.GetApiQueryAuth{}
	err = res.Body.FromJsonTo(GetApiQueryAuthResp)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	err = db.SaveAppidToStore(req.StoreId, GetApiQueryAuthResp.AuthorizationInfo.AuthorizerAppid)
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
