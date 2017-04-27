package service

import (
	"fmt"

	"github.com/franela/goreq"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/weixin/db"
	"github.com/wothing/log"
)

func getandSaveAuthorizerAccountInfo(access_token, component_appid, authorizer_appid string) error {
	type req struct {
		Component_appid  string `json:"component_appid"`
		Authorizer_appid string `json:"authorizer_appid"`
	}

	req_uri := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=%s", access_token)

	log.Debug(req_uri)

	res, err := goreq.Request{
		Method: "POST",
		Uri:    req_uri,
		Body: &req{
			Component_appid:  component_appid,
			Authorizer_appid: authorizer_appid,
		},
	}.Do()
	if err != nil {
		log.Error(err)
		return err
	}
	// str, _ := res.Body.ToString()
	// log.Debug(str)
	// 组织数据结构
	verify_type_info := &pb.VerifyTypeInfo{}
	service_type_info := &pb.ServiceTypeInfo{}
	authorizerInfo := &pb.AuthorizerInfo{
		ServiceTypeInfo: service_type_info,
		VerifyTypeInfo:  verify_type_info,
	}
	authorization_info := &pb.AuthorizationInfo{}
	callback := &pb.GetAuthBaseInfoResp{
		AuthorizerInfo:    authorizerInfo,
		AuthorizationInfo: authorization_info,
	}

	// 解析数据
	err = res.Body.FromJsonTo(callback)
	if err != nil {
		log.Error(err)
		return err
	}
	log.JSONIndent(callback)

	// insert into db
	err = db.SaveAccount(callback)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
