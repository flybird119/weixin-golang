package service

import (
	"fmt"
	"mime"
	"testing"

	"github.com/goushuyun/weixin-golang/pb"
)

func TestUploadFromWeixin(t *testing.T) {
	url := `https://api.weixin.qq.com/cgi-bin/media/get?access_token=jCiexjRIyClaBz_t2TG_SUh7TgmQqc7e-s46qp23GxbGv1QH1R0XEV95HoOI4zGu3yVq16WMQ56BnzY4BiHRNrqB9qhiBaa_UiE8msfpmG8RlzDJn_C7cPEXJ5ZuRN_6PICjAJDTMG&media_id=qsTFkV76Ob-RrUfGvNxISG0vswmZvv5ssRBjy85AevBiwE-j7E5FDsZbzhGvgeHH`

	key, err := upload(url, pb.MediaZone_Test, "12110007027")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(key)

}

func TestUploadLocal(t *testing.T) {
	txt := `attachment; filename="qsTFkV76Ob-RrUfGvNxISG0vswmZvv5ssRBjy85AevBiwE-j7E5FDsZbzhGvgeHH.jpg"`
	mediatype, params, err := mime.ParseMediaType(txt)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("-------%v---------", mediatype)
	fmt.Printf("=========%v=======", params)

}

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
