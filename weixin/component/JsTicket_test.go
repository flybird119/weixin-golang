package component

import (
	"testing"

	"github.com/wothing/log"
)

// func TestGetAccountInfo(t *testing.T) {
// 	service.GetandSaveAuthorizerAccountInfo("eebxHSnIFylljYPieCwi-Lx8BY4gjkOzPXv3aPlVTvQjnJ-hvKib265SRTj9dobfCpiKicfcpvfRZiAhV3iXChxh6SZWsTXR6HSc2OwqyweP75wPX1xd0IQxgXJi33XQXAMgAAAEHT", "wx1c2695469ae47724", "wx6d36779ce4dd3dfa")
// }

func TestJsticket(t *testing.T) {
	ticket, err := JsTicket("wx6d36779ce4dd3dfa", "refreshtoken@@@11TBbayf4_7wqFaWNXEg8Db6EoyuM8m2jIBw-AiI5J4")
	//
	if err != nil {
		log.Error(err)
	}
	t.Log(ticket)
}
