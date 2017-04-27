package component

import (
	"time"
	com "wechat_component/component"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/errs"
	"github.com/wothing/log"

	"github.com/coreos/etcd/client"
	"github.com/goushuyun/weixin-golang/weixin/config"
)

var (
	component              com.WechatComponent = nil
	component_access_token                     = ""
	pre_auth_code                              = ""
)

func init() {
	conf := config.GetConf()
	component = com.New(conf.AppID, conf.AppSecret, conf.AESKey, conf.Token)
}

func ApiAuthorizerToken(appid, refresh_token string) (string, error) {
	// get from etcd
	key := "/bookcloud/weixin/component/AuthorizerAccessToken/" + appid
	resp, err := db.GetEtcdConn().Get(context.Background(), key, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			/*
				token not found at etcd
			*/
			access_token, err := ComponentAccessToken()
			if err != nil {
				log.Error(err)
				return "", err
			}
			publicToken, err := component.GetNormalApi().GetAuthAccessToken(access_token, appid, refresh_token)
			if err != nil {
				log.Error(err)
				return "", err
			}

			// save authorizer token to etcd
			if publicToken.AccessToken != "" {
				_, err = db.GetEtcdConn().Set(context.Background(), key, publicToken.AccessToken, &client.SetOptions{TTL: time.Minute * 90})
				if err != nil {
					return "", errs.NewError(errs.ErrInternal, "etcd error %v", err)
				}
			}

			return publicToken.AccessToken, nil
		} else {
			// other error
			log.Error(err)
			return "", err
		}
	}
	return resp.Node.Value, nil
}

func PreAuthCode() (string, error) {
	access_token, err := ComponentAccessToken()
	if err != nil {
		log.Error(err)
		return "", err
	}
	code, _ := component.GetRegularApi().GetPreAuthCode(access_token)
	return code, nil
}

func ComponentAccessToken() (string, error) {
	// 从etcd中获取 compoment_access_token
	resp, err := db.GetEtcdConn().Get(context.Background(), "/bookcloud/weixin/component/access_token", nil)
	if err != nil {
		if client.IsKeyNotFound(err) {

			log.Debug("access_token not found")

			ticket := Ticket()
			token, _ := component.GetRegularApi().GetAccessToken(ticket)
			if token != "" {
				_, err = db.GetEtcdConn().Set(context.Background(), "/bookcloud/weixin/component/access_token", token, &client.SetOptions{TTL: time.Minute * 90})
				if err != nil {
					return "", errs.NewError(errs.ErrInternal, "etcd error %v", err)
				}
			}
			return token, nil
		} else {
			return "", errs.NewError(errs.ErrInternal, "etcd error %v", err)
		}
	}
	return resp.Node.Value, nil
}
