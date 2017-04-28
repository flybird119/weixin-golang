package component

import (
	"testing"

	"github.com/goushuyun/weixin-golang/weixin/service"
)

func TestGetAccountInfo(t *testing.T) {
	service.GetandSaveAuthorizerAccountInfo("eebxHSnIFylljYPieCwi-Lx8BY4gjkOzPXv3aPlVTvQjnJ-hvKib265SRTj9dobfCpiKicfcpvfRZiAhV3iXChxh6SZWsTXR6HSc2OwqyweP75wPX1xd0IQxgXJi33XQXAMgAAAEHT", "wx1c2695469ae47724", "wx6d36779ce4dd3dfa")
}

func TestJsticket(t *testing.T) {
	ticket, err := JsTicket()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ticket)
}
