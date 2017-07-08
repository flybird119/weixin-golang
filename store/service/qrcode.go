package service

import (
	"fmt"
	"time"

	"github.com/goushuyun/weixin-golang/misc"

	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

func GenQrcode(content string, w, h int64) (string, error) {

	to_fetch_url := "https://pan.baidu.com/share/qrcode?w=%d&h=%d&url=%s"
	to_fetch_url = fmt.Sprintf(to_fetch_url, w, h, content)

	// fetch it and upload to qiniu
	qrcode_key := fmt.Sprintf("recycling_qrcode/%d.png", time.Now().UnixNano())
	fetchImageReq := &pb.FetchImageReq{Zone: 0, Url: to_fetch_url, Key: qrcode_key}

	log.Debug(fetchImageReq)

	resp := &pb.FetchImageResp{}
	err := misc.CallSVC(context.Background(), "bc_mediastore", "FetchImage", fetchImageReq, resp)

	if err != nil {
		log.Debug(err)
		return "", err
	}

	return resp.QiniuUrl, nil
}
