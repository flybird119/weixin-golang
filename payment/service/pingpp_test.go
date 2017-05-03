package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
	"github.com/pingplusplus/pingpp-go/pingpp/refund"
	"github.com/wothing/log"
)

func init() {
	pingpp.LogLevel = 2
	pingpp.Key = "sk_test_ibbTe5jLGCi5rzfH4OqPW9KC"
	fmt.Println("Go SDK Version:", pingpp.Version())
	pingpp.AcceptLanguage = "zh-CN"
	//设置商户的私钥 记得在Ping++上配置公钥
	//pingpp.AccountPrivateKey
}

func TestRefund(t *testing.T) {
	params := &pingpp.RefundParams{
		Amount:      1, //可以注释不上传
		Description: "12345",
	}
	re, err := refund.New("ch_GuPKi1mjXjzDPmz1uD1aPq90", params) //ch_id 是已付款的订单号

	if err != nil {
		log.Fatal(err)
	}
	log.JSONIndent(re)
}

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
	pingpp.AccountPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQC23/9KS+0uVJUGCW/ZFkaCOcBDoSWzQVD3wpOUyqOiKEk9Mpin
fDRvJRKVMecHjMEThA503iWC+TuzcPITbyiSc94ZqNeowzRrKLMqgpXp8xf/iim2
lK8uoz+iwSDA0TR96CYwgeluxgVSLTFQ8E2CD/J4uu61xc7647fVGIqKUwIDAQAB
AoGAAZAIbmoXrL2sSFDsU76M+6/ipLFL0SxNtNBE0pCotUoC1jMIeuXkzM5USlIS
102smK4YMYd0apoWmIHuj5vzjMkdiJqjN0jZIWGu7+QtYyaMbwimBQMfOf5eyofj
P50GEG7aoWwi0cl5Sxmt9zM6DakJmbC8eNNIAOCgypwX/wECQQDfrVDqjZkysuLB
ROvGKdeyu3zSDmYRPkoKaFcwjLXAp6CUMMghHgcppB18PjLAeNzRE4UHjOaDHoYx
bBPCxuXJAkEA0U1BWczW2XWKxl7lindZmFKk96AKgYDXESr/XmS6e6tIyz5smV6j
0dQktB6Z/A7BFvr2Sio9Fmy55LulgHVtOwJAQmnQo8QlX7NTtrUDGJSl8fDPUANs
dOQ80bhHYyf0c16SRE3zrjmfQNL02kYRhaqdTgrwrdw9OWNfzt7bQzMRWQJAWN3G
a4xvhLFFlOhh6aK3JdehN4p6K3Y62o05FCkMjMmzBKiij5QBVmwOkXOUydKx5UH1
JJQ+j7DmVNnfcWVqVQJANmmsVdUjrVR97koRQhGnKjHq93fSC3PWNFD9bFssdO9S
PP56jDrpttNbxDOpYO7ufMLQYNNQhbAo1b+txVFsKQ==
-----END RSA PRIVATE KEY-----`
	metadata := make(map[string]interface{})
	extra := make(map[string]interface{})
	extra["open_id"] = "oWg4qwSn0d3UM0kjZULRdb4SC2hw"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	orderno := r.Intn(999999999999999)

	params := &pingpp.ChargeParams{
		Order_no:  strconv.Itoa(orderno),
		App:       pingpp.App{Id: "app_4qnjLOWXbDKSPmbb"},
		Amount:    1000,
		Channel:   "wx_pub",
		Currency:  "cny",
		Client_ip: "127.0.0.1",
		Subject:   "Your Subject",
		Body:      "Your Body",
		Extra:     extra,
		Metadata:  metadata,
	}

	fmt.Printf(">>>>>>>>>>>>>>>>>>>>>>>>\n%+v\n", params)

	//返回的第一个参数是 charge 对象，你需要将其转换成 json 给客户端，或者客户端接收后转换。
	ch, err := charge.New(params)
	if err != nil {
		errs, _ := json.Marshal(err)
		fmt.Println(string(errs))
		log.Fatal(err)
		return
	}

	log.JSONIndent(ch)

}
