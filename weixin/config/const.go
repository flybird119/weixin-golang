package config

type Conf struct {
	AppID     string
	AppSecret string
	Token     string
	AESKey    string
}

var conf *Conf

func GetConf() *Conf {
	if conf == nil {
		return &Conf{
			AppID:     "wx1c2695469ae47724",
			AppSecret: "ea187a8170f5dc556992d9fa90f6dd56",
			Token:     "goushuyun",
			AESKey:    "goushuyungoushuyungoushuyungoushuyungoushuy",
		}
	}

	return conf
}
