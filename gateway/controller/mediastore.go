package controller

import (
	"net/http"

	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
)

func GetUplaodToken(w http.ResponseWriter, r *http.Request) {
	req := &pb.UpLoadReq{}
	misc.CallWithResp(w, r, "bc_mediastore", "GetUpToken", req, "key")
}
