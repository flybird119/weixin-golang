package component

import (
	"testing"

	"github.com/wothing/log"
)

// func TestGetAccountInfo(t *testing.T) {
// 	service.GetandSaveAuthorizerAccountInfo("eebxHSnIFylljYPieCwi-Lx8BY4gjkOzPXv3aPlVTvQjnJ-hvKib265SRTj9dobfCpiKicfcpvfRZiAhV3iXChxh6SZWsTXR6HSc2OwqyweP75wPX1xd0IQxgXJi33XQXAMgAAAEHT", "wx1c2695469ae47724", "wx6d36779ce4dd3dfa")
// }

func TestJsticket(t *testing.T) {
	ticket, err := JsTicket("wx6d36779ce4dd3dfa", "refreshtoken@@@o509G_ob9H_tvBVV7P2Ba7RbhRmobn1ZV_1GDWgOpOw")
	if err != nil {
		log.Error(err)
	}
	t.Log(ticket)
}
