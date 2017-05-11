package service

import (
	"errors"

	"github.com/goushuyun/weixin-golang/errs"
	"github.com/goushuyun/weixin-golang/misc"
	"golang.org/x/net/context"

	"github.com/goushuyun/weixin-golang/pb"
	"github.com/goushuyun/weixin-golang/user/db"
	"github.com/wothing/log"
)

type UserService struct {
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetUserInfo", "%#v", req))

	// 根据 code , appid（官方）换取用户官方 openid
	weixin_info := &pb.WeixinInfo{}
	err := misc.CallSVC(ctx, "bc_weixin", "GetOfficialOpenid", req, weixin_info)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	// after get official_openid, get user info or save user info
	exist, err := db.OfficalOpenidExist(weixin_info.Openid)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	user := &pb.User{
		WeixinInfo: weixin_info,
	}

	if exist {
		// 获取用户信息
		err := db.GetUserInfoByOfficialOpenid(user)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	} else {
		// save official_openid
		err := db.SaveOfficialOpenid(user)
		if err != nil {
			log.Error(err)
			return nil, errs.Wrap(errors.New(err.Error()))
		}
	}

	return &pb.GetUserInfoResp{}, nil
}

func (s *UserService) GetUserInfoByOpenid(ctx context.Context, req *pb.User) (*pb.User, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "GetUserInfoByOpenid", "%#v", req))

	// get user info by openid
	err := db.GetUserInfo(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return req, nil
}

func (s *UserService) SaveUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	tid := misc.GetTidFromContext(ctx)
	defer log.TraceOut(log.TraceIn(tid, "SaveUser", "%#v", req))

	// save user
	err := db.SaveUser(req)
	if err != nil {
		log.Error(err)
		return nil, errs.Wrap(errors.New(err.Error()))
	}

	return req, nil
}
