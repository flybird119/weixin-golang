package service

import (
	"errors"
	"fmt"
	"gsb/misc"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/weixin/component"
	"github.com/goushuyun/weixin-golang/weixin/config"
	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type WeixinServer struct{}

func (s *WeixinServer) GetAuthURL(ctx context.Context, req *pb.WeixinReq) (*pb.GetAuthURLResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetAuthUrl", "%#v", req))

	conf := config.GetConf()
	pre_auth_code, err := component.PreAuthCode()
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	// get auth url
	redirect_uri := "http://baidu.com"
	url := "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s"

	url = fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s", conf.AppID, pre_auth_code, redirect_uri)

	return &pb.GetAuthURLResp{Code: errs.Ok, Message: "ok", Url: url}, nil
}
