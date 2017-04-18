package db

import (
	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

func SaveUser(user *pb.User) error {
	query := "insert into users(openid, nickname, sex, avatar, store_id) values($1, $2, $3, $4, $5) returing id"

	err := DB.QueryRow(query, user.WeixinInfo.Openid, user.WeixinInfo.Nickname, user.WeixinInfo.Sex, user.WeixinInfo.Headimgurl, user.StoreId).Scan(&user.UserId)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
