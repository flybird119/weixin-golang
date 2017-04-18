package component

import (
	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/log"
)

var ticket = ""

func Ticket() string {
	if ticket == "" {
		resp, err := db.GetEtcdConn().Get(context.Background(), "/bookcloud/weixin/component/ComponentVerifyTicket", nil)

		if err != nil {
			log.Error(err)
			return ""
		}

		if err == nil && resp.Node != nil && resp.Node.Value != "" {
			ticket = resp.Node.Value
			return ticket
		}
	}
	return ticket
}
