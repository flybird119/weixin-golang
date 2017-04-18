package db

import (
	. "github.com/goushuyun/weixin-golang/db"

	"github.com/wothing/log"
)

func SaveAppidToStore(store_id, app_id string) error {
	query := "update store set appid = $1 where id = $2"

	_, err := DB.Exec(query, app_id, store_id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
