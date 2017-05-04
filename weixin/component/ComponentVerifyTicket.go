package component

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/log"
)

func Ticket() (string, error) {
	// 每次从 etcd 取 ticket, 若etcd中没有，则抛错

	resp, err := db.GetEtcdConn().Get(context.Background(), "/bookcloud/weixin/component/ComponentVerifyTicket", nil)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if err == nil && resp.Node != nil && resp.Node.Value != "" {
		return resp.Node.Value, nil
	}

	return "", errors.New("there is no ticket in etcd")
}
