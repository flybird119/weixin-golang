package service

import (
	"testing"

	"github.com/goushuyun/weixin-golang/pb"
)

func TestFetch(t *testing.T) {

	t.Log("My name is Wang Kai ...")

	url, err := FetchImg(pb.MediaZone_Public, "http://wx4.sinaimg.cn/mw690/5a3cbcf7ly1fbo9k0kiukj21ve2io7wo.jpg", "beauty")

	t.Log(err, url)
}

func TestMakeToken(t *testing.T) {
	token, url := makeToken(0, "wanghaiting")

	t.Log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	t.Log(token, url)
	t.Log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
}

func TestGenAt(t *testing.T) {
	t.Log(GenAccessToken("/v2/tune/refresh\n"))
}

func TestRefreshUrls(t *testing.T) {
	err := RefreshURLCache([]string{"http://image.cumpusbox.com/book/9787513557344"})

	if err != nil {
		t.Log(err)
	}
}
