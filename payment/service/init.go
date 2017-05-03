package service

import (
	"fmt"

	"github.com/pingplusplus/pingpp-go/pingpp"
)

func init() {
	// LogLevel 是 Go SDK 提供的 debug 开关
	pingpp.LogLevel = 2
	//设置 API Key
	pingpp.Key = "sk_live_mHebTSOm1S0G8y5SW5zTSaDO"
	//获取 SDK 版本
	fmt.Println("Go SDK Version:", pingpp.Version())
	//设置错误信息语言，默认是中文
	pingpp.AcceptLanguage = "zh-CN"

	//设置商户的私钥 记得在Ping++上配置公钥
	pingpp.AccountPrivateKey = privity_key
}

var privity_key = `-----BEGIN RSA PRIVATE KEY-----
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

var public_key = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC23/9KS+0uVJUGCW/ZFkaCOcBD
oSWzQVD3wpOUyqOiKEk9MpinfDRvJRKVMecHjMEThA503iWC+TuzcPITbyiSc94Z
qNeowzRrKLMqgpXp8xf/iim2lK8uoz+iwSDA0TR96CYwgeluxgVSLTFQ8E2CD/J4
uu61xc7647fVGIqKUwIDAQAB
-----END PUBLIC KEY-----`
