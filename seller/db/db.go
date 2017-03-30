package db

import (
	"database/sql"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	log "qiniupkg.com/x/log.v7"
)

func Save() error {
	timestamp := time.Now().Unix()
	query := "insert into seller (mobile,password,username,name,avatar,register_time,status,id_card) values($1,$2,$3,$4,$5,$6,$7,$8)"
	_, err := DB.Exec(query, "13122210065", "492226568", "orican", "李肖", "image", timestamp, 0, "411627199107295410")
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//CheckSellerExists 通过手机号和登录密码检查商家是否存在
func CheckSellerExists(loginModel *pb.LoginModel) (userinfo *pb.UserInfo, err error) {
	query := "select * from seller as s where s.mobile=$1 and s.password=$2"
	log.Debugf("select id,mobile,username from seller as s where s.mobile=%s and s.password=%s", loginModel.Mobile, loginModel.Password)
	//如果检查失败
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		log.Error(err)
		return nil, err
	default:
		return userinfo, nil
	}

	//SellerRegister 商家注册
	func SellerRegister(registerModel *pb.RegisterModel) error {
		query := "insert into seller (mobile,password,username,create_at,update_at,status,)"
		fmt.Println(query)
		return nil
	}

	//private functions
	func encryPassword(password string) (string) {
		h := md5.New()
		h.Write([]byte(password)) // 需要加密的字符串为 sharejs.com
		result := hex.EncodeToString(h.Sum(nil))
		return result;
	}
}
