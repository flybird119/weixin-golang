package component

import (
	"fmt"
	"goushuyun/errs"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc/http"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type ticket_callback struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Ticket  string `json:"ticket"`
}

// 接口调用次数有限，需存入etcd
func JsTicket(appid, refresh_token string) (string, error) {
	key := "/bookcloud/weixin/component/js_ticket/" + appid
	resp, err := db.GetEtcdConn().Get(context.Background(), key, nil)
	if err != nil {
		if client.IsKeyNotFound(err) {
			/*
				js_ticket not found
			*/
			// get access_token
			access_token, err := ApiAuthorizerToken(appid, refresh_token)
			if err != nil {
				log.Error(err)
				return "", err
			}

			// get request to get ticket
			callback := &ticket_callback{}

			url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", access_token)
			err = http.GETWithUnmarshal(url, callback)
			if err != nil {
				log.Error(err)
				return "", err
			}

			// save js_ticket to etcd
			if callback.Ticket != "" {
				_, err = db.GetEtcdConn().Set(context.Background(), key, callback.Ticket, &client.SetOptions{TTL: time.Minute * 100})
				if err != nil {
					return "", errs.NewError(errs.ErrInternal, "etcd error %v", err)
				}
			}

			return callback.Ticket, nil

		} else {
			log.Error(err)
			return "", err
		}
	}
	return resp.Node.Value, nil

}
