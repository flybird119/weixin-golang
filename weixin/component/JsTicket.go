package component

import (
	"fmt"

	"github.com/goushuyun/weixin-golang/misc/http"
	"github.com/wothing/log"
)

type ticket_callback struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Ticket  string `json:"ticket"`
}

func JsTicket(appid, refresh_token string) (string, error) {
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

	return callback.Ticket, nil
}
