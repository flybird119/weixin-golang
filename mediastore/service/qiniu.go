package service

import (
	"goushuyun/pb"

	"github.com/wothing/log"
	"qiniupkg.com/api.v7/kodo"
	"qiniupkg.com/api.v7/kodocli"
)

var cli *kodo.Client
var uploader kodocli.Uploader

const (
	ak = "PRBrdNUf1m07enHI26CgpG7z3O8YLYXgoCE-1P2S"
	sk = "PtCRtf7xyY0P6mBaQARL6KeYW6dVdglFzIz3BiQS"
)

type qiniuZone struct {
	public bool
	bucket string
	domain string
	imgsuf string
	thmsuf string
}

var zones = map[string]*qiniuZone{
	"Test": &qiniuZone{
		public: true,
		bucket: "goushuyun-test",
		domain: "http://ojrfwndal.bkt.clouddn.com/",
		imgsuf: "",
		thmsuf: "-th",
	},
	"Public": &qiniuZone{
		public: true,
		bucket: "goushuyun-public",
		domain: "http://image.cumpusbox.com/",
		imgsuf: "",
		thmsuf: "-th",
	},
	"Goushubao": &qiniuZone{
		public: true,
		bucket: "gbs-public",
		domain: "http://olto51num.bkt.clouddn.com/",
		imgsuf: "",
		thmsuf: "-th",
	},
}

func init() {
	kodo.SetMac(ak, sk)
	zone := 0
	cli = kodo.New(zone, nil)
}

func makeToken(zone pb.MediaZone, filename string) (token, url string) {
	var c *qiniuZone
	if zone <= pb.MediaZone_Goushubao && zone >= pb.MediaZone_Test {
		c = zones[zone.String()]
	} else {
		c = zones["Test"]
	}
	scope := c.bucket

	if filename != "" {
		scope = scope + ":" + filename
	}

	log.Debugf("The scope is : ", scope)

	policy := &kodo.PutPolicy{
		Scope:   scope,
		Expires: 600,
	}

	token = cli.MakeUptoken(policy)
	url = c.domain + filename

	return
}

func FetchImg(zone pb.MediaZone, url, key string) (string, error) {
	z := zones[zone.String()]

	log.Debugf("The zone is : %+v", z)

	p := cli.Bucket(z.bucket)
	err := p.Fetch(nil, key, url)
	if err != nil {
		return "", err
	}
	log.Debugf("Fetching successful, The url is: ", z.domain+key)
	return z.domain + key, nil
}
