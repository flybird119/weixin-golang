package db

import (
	"database/sql"
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

func SaveOfficialOpenid(user *pb.User) error {
	query := "insert into users(official_openid) values($1) returning id"
	err := DB.QueryRow(query, user.WeixinInfo.Openid).Scan(&user.UserId)

	log.Debugf("insert into users(official_openid) values('%s') returning id", user.WeixinInfo.Openid)

	if err != nil {
		log.Error(err)
		return nil
	}
	return nil
}

func getMyStoreOpenid(user_id, store_id string) (string, error) {
	// 获取某个用户 在某个云店铺的 store_id
	var openid string
	query := "select openid from map_store_users where user_id = $1 and store_id = $2"

	err := DB.QueryRow(query, user_id, store_id).Scan(&openid)

	if err == sql.ErrNoRows {
		log.Debug("The user has no openid in this store")
		return "", nil
	}

	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("select openid from map_store_users where user_id = '%s' and store_id = '%s'", user_id, store_id)

	return openid, nil
}

func GetUserInfoByOfficialOpenid(user *pb.User) error {
	// 备注：map_store_users 中的值是在获取用户微信信息时写入的，所以在这里需要使用左连接SQL（以为此时值不一定有）
	// 取出用户的基本信息、对应store_id 的 openid
	query := "select id, nickname, sex, avatar, status from users where official_openid = $1"

	err := DB.QueryRow(query, user.WeixinInfo.Openid).Scan(&user.UserId, &user.WeixinInfo.Nickname, &user.WeixinInfo.Sex, &user.WeixinInfo.Headimgurl, &user.Status)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debugf("select id, nickname, sex, avatar, status from users where official_openid = '%s'", user.WeixinInfo.Openid)

	// 获取用户对应店铺的 openid
	currentStoreOpenid, err := getMyStoreOpenid(user.UserId, user.StoreId)
	if err != nil {
		log.Error(err)
		return err
	}

	user.CurrentStoreOpenid = currentStoreOpenid
	return nil
}

func OfficalOpenidExist(official_openid string) (bool, error) {
	var total int64
	query := "select count(*) from users where official_openid = $1"
	err := DB.QueryRow(query, official_openid).Scan(&total)
	log.Debugf("select count(*) from users where official_openid = '%s'", official_openid)

	if err != nil {
		log.Error(err)
		return false, err
	}

	return total > 0, nil
}

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
		// insert user to DB
		err = insertUser(user)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func insertUser(user *pb.User) error {
	store_ids := fmt.Sprintf("{\"%s\"}", user.StoreId)

	log.Debug(">>>>>>>>>>>>>", store_ids, "<<<<<<<<<<<<<<<<")

	query := "insert into users(openid, nickname, sex, avatar, store_ids) values('%s', '%s', %d, '%s', '%s') returning id"
	err := DB.QueryRow(fmt.Sprintf(query, user.WeixinInfo.Openid, user.WeixinInfo.Nickname, user.WeixinInfo.Sex, user.WeixinInfo.Headimgurl, store_ids)).Scan(&user.UserId)
	log.Debug(fmt.Sprintf(query, user.WeixinInfo.Openid, user.WeixinInfo.Nickname, user.WeixinInfo.Sex, user.WeixinInfo.Headimgurl, store_ids))

	if err != nil {
		log.Error(err)
		return err
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
