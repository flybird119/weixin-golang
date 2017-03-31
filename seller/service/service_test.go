package service

import (
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/log"
)

func TestRedis(t *testing.T) {
	db.InitRedis("127.0.0.1")
	var con redis.Conn = db.GetRedisConn()
	defer con.Close()
	_, err := con.Do("set", "seller", "nihao")
	if err != nil {
		fmt.Println("111")
		log.Debug(err)
		return
	}
	_, err = con.Do("expire", "seller", 5)
	if err != nil {
		fmt.Println("222")
		log.Debug(err)
		return
	}
	resultStr, err := redis.String(con.Do("get", "seller11111"))
	if err == redis.ErrNil {
		log.Debug("not found")
		return
	}
	if err != nil {

		log.Debug(err)
		return
	}
	log.Debug(">>>>>%s", resultStr)

	err = con.Flush()
	if err != nil {
		fmt.Println("444")
		log.Debug(err)
		return
	}
	err = con.Close()

	if err != nil {
		fmt.Println("555")
		log.Debug(err)
		return
	}
}
