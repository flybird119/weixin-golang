package misc

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/wothing/log"
)

func TestPgArray2StringSlice(t *testing.T) {
	cases := []struct {
		str      []byte
		expected []string
	}{
		{[]byte("{abc,de,f}"), []string{"abc", "de", "f"}},
		{[]byte("{abc}"), []string{"abc"}},
		{[]byte("{abc,}"), []string{"abc"}},
		{[]byte("{abc} "), []string{"abc"}},
		{[]byte("{abc"), []string{"abc"}},
		{[]byte("{}"), []string{}},
		{[]byte("}"), []string{}},
		{[]byte("{"), []string{}},
		{[]byte("x"), []string{"x"}},
		{nil, []string{}},
	}
	for i, c := range cases {
		actual := PgArray2StringSlice(c.str)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("case %d expected string array %#v, but got %#v", i, c.expected, actual)
		}
	}
}

func TestTransaction(t *testing.T) {
	db.InitPG("hello")
	defer db.ClosePG()
	for i := 0; i < 10; i++ {
		go update1("mobile"+strconv.Itoa(i), "nickname"+strconv.Itoa(i), "00000001")
	}
	time.Sleep(15 * time.Second)
}

func update1(mobile, nickname, id string) {
	tx, _ := db.DB.Begin()
	var status int

	_, err := tx.Exec("update seller set mobile=$1 ,nickname=$2,status=status+1 where id=$3", mobile, nickname, id)

	if err != nil {
		log.Debugf("执行出错:%+v", err)

		return
	}
	log.Debugf("查询开始")
	err = tx.QueryRow("select status from seller where id=$1", id).Scan(&status)
	if err != nil {
		log.Debugf("查询出错:%+v", err)
		return
	}
	log.Debugf("查询结束")
	log.Debugf("开始状态，%d", status)
	if status > 30 {
		log.Debugf("事务回滚")
		err = tx.Rollback()
		if err != nil {
			log.Debugf("提交出错%+v", err)
		}
		return
	}
	log.Debugf("事务提交开始")
	err = tx.Commit()
	log.Debugf("事务提交结束")
	if err != nil {
		log.Debugf("提交出错%+v", err)
		return
	}
}
