package db

import (
	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

func SaveUser(user *pb.User) error {
	// check if user is exist
	isExist, err := isExist(user)
	if err != nil {
		log.Error(err)
		return nil
	}

	if isExist {
		err = GetUserInfo(user)
		if err != nil {
			log.Error(err)
			return nil
		}
	} else {
		query := "insert into users(openid, nickname, sex, avatar, store_id) values($1, $2, $3, $4, $5) returning id"
		log.Debugf("insert into users(openid, nickname, sex, avatar, store_id) values('%s', '%s', %d, '%s', '%s') returning id", user.WeixinInfo.Openid, user.WeixinInfo.Nickname, user.WeixinInfo.Sex, user.WeixinInfo.Headimgurl, user.StoreId)

		err = DB.QueryRow(query, user.WeixinInfo.Openid, user.WeixinInfo.Nickname, user.WeixinInfo.Sex, user.WeixinInfo.Headimgurl, user.StoreId).Scan(&user.UserId)

		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func isExist(user *pb.User) (bool, error) {
	query := "select count(*) from users where openid = $1"
	var total int64
	log.Debugf("select count(*) from users where openid = '%s'", user.WeixinInfo.Openid)
	err := DB.QueryRow(query, user.WeixinInfo.Openid).Scan(&total)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return total > 0, nil
}

func GetUserInfo(user *pb.User) error {
	query := "select id, nickname, sex, avatar, status, store_id from users where openid = $1"
	log.Debugf("select id, nickname, sex, avatar, status, store_id from users where openid = '%s'", user.WeixinInfo.Openid)

	DB.QueryRow(query, user.WeixinInfo.Openid).Scan(&user.UserId, &user.WeixinInfo.Nickname, &user.WeixinInfo.Sex, &user.WeixinInfo.Headimgurl, &user.Status, &user.StoreId)

	return nil
}
