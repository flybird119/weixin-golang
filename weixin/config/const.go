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
			AppSecret: "bd8f125dc0300451e7495c70f5480575",
			Token:     "goushuyun",
			AESKey:    "goushuyungoushuyungoushuyungoushuyungoushuy",
		}
	}

	return conf
}
