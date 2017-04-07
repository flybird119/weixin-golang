package service

import (
	"errors"
	"goushuyun/errs"

	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
	"golang.org/x/net/context"
)

type MediastoreServer struct {
	Test bool
}

func (s *MediastoreServer) RefreshUrls(ctx context.Context, req *pb.RefreshReq) (*pb.NormalResp, error) {
	// to refresh urls
	err := RefreshURLCache(req.Urls)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return &pb.NormalResp{Code: errs.Ok, Message: "ok"}, nil
}

func (s *MediastoreServer) GetUpToken(ctx context.Context, req *pb.UpLoadReq) (*pb.GetUpTokenResp, error) {
	token, url := makeToken(req.Zone, req.Key)

	return &pb.GetUpTokenResp{Token: token, Url: url}, nil
}

func (s *MediastoreServer) FetchImage(ctx context.Context, req *pb.FetchImageReq) (*pb.FetchImageResp, error) {
	url, err := FetchImg(req.Zone, req.Url, req.Key)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}
	return &pb.FetchImageResp{QiniuUrl: url}, nil
}
