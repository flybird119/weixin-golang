package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"

	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
)

func TestPay(t *testing.T) {
	// LogLevel 是 Go SDK 提供的 debug 开关
	pingpp.LogLevel = 2
	//设置 API Key
	pingpp.Key = "sk_test_v94yzLeL8S8COiTGC40yLiDK"
	//获取 SDK 版本
	fmt.Println("Go SDK Version:", pingpp.Version())
	//设置错误信息语言，默认是中文
	pingpp.AcceptLanguage = "zh-CN"

	//设置商户的私钥 记得在Ping++上配置公钥
	pingpp.AccountPrivateKey = `
    -----BEGIN RSA PRIVATE KEY-----
    MIICXQIBAAKBgQDeyIgIuz+bRZHtcGxKIvJ/srOMOUcKDXE0ek4c/HOtSzdcWxcw
    28gTVfy+W4Be6+Ix17xyKIxiQiSMF5xYRvNRa1oK8UiPTUnQHKZRDgRwPjCaZk5X
    fBV0wvYmbCZbfb0xylpvb0it7Bwt1YdxqtFboR56M46S8n0U3qKnGn5TcQIDAQAB
    AoGBAMrcFuK8fqLIqqRmpnSrdd1Jv6yDy2gf7WE3rUE/r6Www+xZFbjrqDfTKJ29
    fBry97kjFPluasZeLCFUrozDrnJVwMT6hewQLvTSyZuQCigZBN8DJVx5yuVE16d3
    OxHYDGdj8C36ZupRfI9DWeWFj+zWzFTCgsd+5mZjhnJkQNJBAkEA+m6ukyiAVXmB
    wrGzb75sySErpom8wCHUSdxLmyF9MrIfIOGCrVyECKqoK/sP7jHvkip1wKvXTeO0
    SjJ0qmSgnQJBAOO8fKog7Cjtio5VIHfdEQREky2jNsjkeqwT0DeANHiXdKyzPb0L
    oXhW+TJJO8zDckaLXVJnsSbC+6y7d34bE+UCQQDCvOK/yBTTYqMG1MwlrrxFQqgA
    3saJ2USNEuMwBMCodV5DYVkOmgyJ+LrBSH/Ax8/1p1LdukK4bMK7l7Sk8475AkAy
    C7/Rmz6Kl/j04lwqOxh8OZ2mT9HAQAV9PzVonPHq9k2bjiApJR8s1OAaXuGXU/QO
    8J1neIYDoKGyCdhujADJAkABRuWXkJms1pT78yPMa0zN4PIOq2R0UMlEl+GpMG7a
    b+Pv1GBCJAkDq/mktvA7t9YVPcPOXkwQzWtVTxN0JOH7
    -----END RSA PRIVATE KEY-----
`
	metadata := make(map[string]interface{})
	metadata["color"] = "red"
	extra := make(map[string]interface{})

	pingpp.Key = "sk_test_v94yzLeL8S8COiTGC40yLiDK"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderno := r.Intn(999999999999999)

	params := &pingpp.ChargeParams{
		Order_no:  strconv.Itoa(orderno),
		App:       pingpp.App{Id: "app_1Gqj58ynP0mHeX1q"},
		Amount:    1000,
		Channel:   "alipay_wap",
		Currency:  "cny",
		Client_ip: "127.0.0.1",
		Subject:   "Your Subject",
		Body:      "Your Body",
		Extra:     extra,
		Metadata:  metadata,
	}

	//返回的第一个参数是 charge 对象，你需要将其转换成 json 给客户端，或者客户端接收后转换。
	ch, err := charge.New(params)
	if err != nil {
		errs, _ := json.Marshal(err)
		fmt.Println(string(errs))
		log.Fatal(err)
		return
	}

	fmt.Println(ch)

}
