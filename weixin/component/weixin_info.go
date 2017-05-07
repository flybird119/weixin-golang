package component

import (
	"fmt"

	"github.com/goushuyun/weixin-golang/misc/http"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/weixin/config"
	"github.com/wothing/log"
)

func GetWeixinInfo(req *pb.WeixinReq) (*pb.WeixinInfo, error) {
	component_access_token, err := ComponentAccessToken()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	config := config.GetConf()

	log.Debug("++++++++++++++++获取微信信息的请求++++++++++++++++++")
	log.JSONIndent(req)
	log.Debugf("--------------component_access_token: %s--------------", component_access_token)

	log.Debug("++++++++++++++++++++++++++++++++++")

	// get access_token
	access_token_url := "https://api.weixin.qq.com/sns/oauth2/component/access_token?appid=%s&code=%s&grant_type=authorization_code&component_appid=%s&component_access_token=%s"

	access_token_url = fmt.Sprintf(access_token_url, req.Appid, req.Code, config.AppID, component_access_token)
	getAcessTokenResp := &GetAcessTokenResp{}
	err = http.GETWithUnmarshal(access_token_url, getAcessTokenResp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Debugf("------------授权 access_token: %s--------------", getAcessTokenResp.AccessToken)

	// get user info
	get_user_info_url := "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	get_user_info_url = fmt.Sprintf(get_user_info_url, getAcessTokenResp.AccessToken, getAcessTokenResp.Openid)
	weixinInfo := &pb.WeixinInfo{}
	err = http.GETWithUnmarshal(get_user_info_url, weixinInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return weixinInfo, nil
}
