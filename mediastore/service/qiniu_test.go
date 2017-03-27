package service

import (
	"goushuyun/pb"
	"testing"
)

func TestFetch(t *testing.T) {

	t.Log("My name is Wang Kai ...")

	url, err := FetchImg(pb.MediaZone_Public, "http://wx4.sinaimg.cn/mw690/5a3cbcf7ly1fbo9k0kiukj21ve2io7wo.jpg", "beauty")

	t.Log(err, url)
}

func TestMakeToken(t *testing.T) {
	token, url := makeToken(0, "haiting")

	t.Log(token, url)
}
