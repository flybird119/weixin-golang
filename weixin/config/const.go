package config

import (
	"github.com/goushuyun/weixin-golang/db"
)

type Conf struct {
	AppID     string
	AppSecret string
	Token     string
	AESKey    string
}

var conf *Conf

const (
	svcName = "weixin"
)

func GetConf() *Conf {
	// get AppID & AppSecret from etcd

	if conf == nil {
		return &Conf{
			AppID:     db.GetValue(svcName, "/component/AppID", "wx1c2695469ae47724"),
			AppSecret: db.GetValue(svcName, "/component/AppSecret", "bd8f125dc0300451e7495c70f5480575"),
			Token:     "goushuyun",
			AESKey:    "goushuyungoushuyungoushuyungoushuyungoushuy",
		}
	}

	return conf
}
