package service

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/robfig/cron"
)

func TestCronPolling(t *testing.T) {
	fmt.Println(os.Getenv("GOPATH"))
	c := cron.New()
	c.AddFunc("0/5 * * * * *", func() { fmt.Println("every five seconds") })
	c.AddFunc("0 0/1 * * * *", func() { fmt.Println("every minutes") })
	//c.AddFunc("0/5 * * * * ?", test) //5秒执行一次，12×5=60，所以一共执行12次
	c.Start()
	//log.Info("Every hour on the half hour")
	defer c.Stop()
	time.Sleep(time.Minute) //一分钟后主线程退出
	fmt.Println("aaa")
}

func TestDate(t *testing.T) {
	tm := time.Now()
	tm = tm.AddDate(0, 0, 14)
	expireDate := tm.Format("2006-01-02日")
	fmt.Println(expireDate)
}
