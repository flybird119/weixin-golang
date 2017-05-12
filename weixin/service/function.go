package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/franela/goreq"
	globalDB "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/weixin/db"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

func getUserBaseInfo(openid, access_token string) (*pb.WeixinInfo, error) {
	url := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"

	url = fmt.Sprintf(url, access_token, openid)

	weixin_info := &pb.WeixinInfo{}
	resp, err := goreq.Request{Method: "POST", Uri: url}.Do()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = resp.Body.FromJsonTo(weixin_info)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return weixin_info, nil
}

func saveAuthorizerAccessTokenToEtcd(appid, token string) error {
	key := "/bookcloud/weixin/component/AuthorizerAccessToken/" + appid
	_, err := globalDB.GetEtcdConn().Set(context.Background(), key, token, &client.SetOptions{TTL: time.Minute * 100})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

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

func DicSort(strs ...string) string {
	sort.Strings(strs)
	return strings.Join(strs, "")
}

func Sha1Str(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetRandomString generate random string by specify chars.
func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(bytes)
}
