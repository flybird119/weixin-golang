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
		/*
			if this user has exist, check the req's store_id whether in DB's store_ids or not
		*/
		storeId_has_exist, err := storeId_has_exist(user.StoreId, user.WeixinInfo.Openid)
		if err != nil {
			log.Error(err)
			return err
		}

		if !storeId_has_exist {
			// 这个user还没有绑定这个 store_id，将新的store_id append 进去
			err = appendNewStoreId(user.StoreId, user.WeixinInfo.Openid)
			if err != nil {
				log.Error(err)
				return nil
			}
		}

		// store_id 已经绑定到了这个 user, 获取用户信息，并返回
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

func appendNewStoreId(store_id, openid string) error {
	// append a new store_id
	query := "update users set store_ids = array_append(store_ids, $1) where openid = $2"
	_, err := DB.Exec(query, store_id, openid)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func storeId_has_exist(store_id, open_id string) (bool, error) {
	query := "select count(*) from users where openid = $1 and $2 = ANY(store_ids)"
	var total int64

	err := DB.QueryRow(query, open_id, store_id).Scan(&total)
	if err != nil {
		log.Error(err)
		return false, err
	}

	return total > 0, nil
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
	query := "select id, nickname, sex, avatar, status from users where openid = $1"
	log.Debugf("select id, nickname, sex, avatar, status from users where openid = '%s'", user.WeixinInfo.Openid)

	DB.QueryRow(query, user.WeixinInfo.Openid).Scan(&user.UserId, &user.WeixinInfo.Nickname, &user.WeixinInfo.Sex, &user.WeixinInfo.Headimgurl, &user.Status)

	return nil
}
