package controller

import (
	"context"
	"encoding/xml"
	"net/http"
	"wechat_component/lib"

	"github.com/goushuyun/weixin-golang/misc/token"

	"github.com/goushuyun/weixin-golang/errs"

	"github.com/coreos/etcd/client"
	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

func WeixinDBTest(w http.ResponseWriter, r *http.Request) {
	req := &pb.WeixinReq{}
	misc.CallWithResp(w, r, "bc_weixin", "TestWeixinDB", req)
}

func GetWeixinInfo(w http.ResponseWriter, r *http.Request) {
	req := &pb.WeixinReq{}
	misc.CallWithResp(w, r, "bc_weixin", "GetWeixinInfo", req, "store_id", "code")
}

func GetApiQueryAuth(w http.ResponseWriter, r *http.Request) {
	req := &pb.WeixinReq{}
	if c := token.Get(r); c != nil {
		req.StoreId = c.StoreId
	} else {
		misc.RespondMessage(w, r, map[string]interface{}{
			"code":    errs.ErrTokenNotFound,
			"message": "token not found",
		})
	}
	misc.CallWithResp(w, r, "bc_weixin", "GetOfficialAccountInfo", req, "auth_code")
}

func GetAuthURL(w http.ResponseWriter, r *http.Request) {
	req := &pb.WeixinReq{}
	misc.CallWithResp(w, r, "bc_weixin", "GetAuthURL", req)
}

const (
	wx_token = "goushuyun"
	aesKey   = "goushuyungoushuyungoushuyungoushuyungoushuy"
	appID    = "wx1c2695469ae47724"
)

var c *lib.MessageCrypter = nil

func getCrypter() (*lib.MessageCrypter, error) {
	if c != nil {
		return c, nil
	}
	crypter, err := lib.NewmessageCrypter(wx_token, aesKey, appID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	c = &crypter
	return c, nil
}

type resp struct {
	XMLName xml.Name `xml:"xml"`
	Encrypt string   `xml:"Encrypt"`
	AppId   string   `xml:"AppId"`
}

type ticketXML struct {
	XMLName               xml.Name `xml:"xml"`
	AppID                 string   `xml:"AppID"`
	CreateTime            string   `xml:"CreateTime"`
	ComponentVerifyTicket string   `xml:"ComponentVerifyTicket"`
}

func ReceiveTicket(w http.ResponseWriter, r *http.Request) {
	log.Debugf("The request body is : %s", r.Context().Value("body"))
	// save ticket into etcd
	v := &resp{}

	err := xml.Unmarshal([]byte(r.Context().Value("body").([]byte)), v)
	if err != nil {
		log.Error(err)
		misc.RespondMessage(w, r, "fail")
	}

	// decrypt
	c, err := getCrypter()
	if err != nil {
		log.Error(err)
		misc.RespondMessage(w, r, "fail")
	}
	crypter := *c
	crypterText, _, err := crypter.Decrypt(v.Encrypt)
	if err != nil {
		log.Error(err)
		misc.RespondMessage(w, r, "fail")
	}

	// get component_verify_ticket, put it into etcd
	ticketXML := &ticketXML{}
	err = xml.Unmarshal(crypterText, ticketXML)
	if err != nil {
		log.Error(err)
		misc.RespondMessage(w, r, "fail")
	}

	_, err = db.GetEtcdConn().Set(context.Background(), "bookcloud/weixin/component/ComponentVerifyTicket", ticketXML.ComponentVerifyTicket, &client.SetOptions{})
	if err != nil {
		log.Error(err)
		misc.RespondMessage(w, r, "fail")
	}

	log.Debugf("the cryter xml is : %s", ticketXML.ComponentVerifyTicket)
}
