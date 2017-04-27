package service

import (
	"fmt"
	"testing"

	"github.com/franela/goreq"
	"github.com/goushuyun/weixin-golang/pb"
)

func TestReq(t *testing.T) {
	uri := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s", "XttAxdfXvyItNvhH5UX4JQSXe3viDrnhvPiEnCMdD7wUu9VYl0sVvjGdvnXn6NWWXjRIPQ7EN1lZZxFyAbFbkwQgavqx2yty_l-Gda0spMcVWQhAAAGHB")

	type Item struct {
		ComponentAppid    string `json:"component_appid"`
		AuthorizationCode string `json:"authorization_code"`
	}

	item := &Item{
		ComponentAppid:    "wx1c2695469ae47724",
		AuthorizationCode: "queryauthcode@@@fSqiFGNkZWnIOpl8JNZVnHA4PVG6jFfuyuaKU4IdXtgxNAMCKYoCrC9AkM260Yge0AmHu8775PrJWdzj4HRQ-A",
	}

	res, err := goreq.Request{
		Method: "POST",
		Uri:    uri,
		Body:   item,
	}.Do()

	if err != nil {
		t.Fatal(err)
	}

	GetApiQueryAuthResp := &pb.GetApiQueryAuth{}
	err = res.Body.FromJsonTo(GetApiQueryAuthResp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.Body.ToString())

	t.Logf("%+v", GetApiQueryAuthResp)
}
