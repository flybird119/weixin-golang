package db

import (
	"time"

	"github.com/wothing/log"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
)

//AddStore 通过手机号和登录密码检查商家是否存在
func AddStore(store *pb.Store) error {
	query := "insert into store (name,status,expire_at) values($1,$2,$3) returning id,extract(epoch from create_at)::integer "
	now := time.Now()
	now = now.Add(7 * 24 * time.Hour)
	err := DB.QueryRow(query, store.Name, pb.StoreStatus_Normal, now).Scan(&store.Id, &store.CreateAt)
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}

func UpdateStore(store *pb.Store) error {
	query := "update store set name=$1,profile=$2"
	_, err := DB.Query(query, store.Name, store.Profile)
	if err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
