package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//CheckSellerExists 通过手机号和登录密码检查商家是否存在
func CheckSellerExists(loginModel *pb.LoginModel) (*pb.SellerInfo, error) {
	query := "select id,mobile,nickname from seller as s where s.mobile=$1 and s.password=$2"
	log.Debugf("select id,mobile,nickname from seller as s where s.mobile=%s and s.password=%s", loginModel.Mobile, loginModel.Password)

	sellerInfo := &pb.SellerInfo{}
	loginModel.Password = encryPassword(loginModel.Password)
	err := DB.QueryRow(query, loginModel.Mobile, loginModel.Password).Scan(&sellerInfo.Id, &sellerInfo.Mobile, &sellerInfo.Username)
	//如果检查失败
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		log.Error(err)
		return nil, err
	default:
		return sellerInfo, nil
	}
}

//SellerRegister 商家注册
func SellerRegister(registerModel *pb.RegisterModel) (id string, err error) {
	//首先检查是否存在相同的手机号好吗
	isExist := CheckMobileExist(registerModel.Mobile)
	if isExist {
		return "", nil
	}
	query := "insert into seller (mobile,password,nickname,status) values (%s) returning id"
	registerModel.Password = encryPassword(registerModel.Password)
	params := fmt.Sprintf("'%s','%s','%s',%d", registerModel.Mobile, registerModel.Password, registerModel.Username, 0)
	query = fmt.Sprintf(query, params)
	log.Debugf("=============>", query)
	err = DB.QueryRow(query).Scan(&id)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return id, nil
}

//CheckMobileExist 检查注册手机号是否存在
func CheckMobileExist(mobile string) bool {
	query := "select id from seller s where s.mobile=$1"
	log.Debugf("select id from seller s where s.mobile=%s", mobile)
	id := ""
	DB.QueryRow(query, mobile).Scan(&id)
	if id == "" {
		return false
	}

	return true

}

//GetSellerByMobile 通过手机号拿到用户的id
func GetSellerByMobile(mobile string) (string, error) {
	query := "select id from seller s where s.mobile=$1"
	log.Debugf("select id from seller s where s.mobile=%s", mobile)
	id := ""
	err := DB.QueryRow(query, mobile).Scan(&id)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return id, nil

}

//GetStoresBySeller 通过商家获取所管理的店铺
func GetStoresBySeller(seller *pb.SellerInfo) (s []*pb.SelfStoresResp_Store, err error) {
	query := "select s.id,s.name,s.logo,extract(epoch from s.expire_at)::integer,extract(epoch from s.create_at)::integer,ms.role from store s  join map_store_seller ms on  s.id=ms.store_id where ms.seller_id=$1 order by id "

	log.Debugf("select s.id,s.name,s.logo,extract(epoch from s.expire_at)::integer,extract(epoch from s.create_at)::integer,ms.role from store s  join map_store_seller ms on  s.id=ms.store_id where ms.seller_id=%s order by id", seller.Id)

	rows, err := DB.Query(query, seller.Id)
	//

	if err != nil {
		return nil, err
	}
	var logo sql.NullString
	defer rows.Close()
	for rows.Next() {
		var store pb.SelfStoresResp_Store
		s = append(s, &store)
		err = rows.Scan(&store.Id, &store.Name, &logo, &store.ExpireAt, &store.CreateAt, &store.Role)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		if logo.Valid {
			store.Logo = logo.String
		} else {
			store.Logo = ""
		}
	}
	if err = rows.Err(); err != nil {
		log.Debug("scan rows err last error: %s", err)
		return nil, err
	}

	return s, nil
}

//UpdatePassword 更改手机号
func UpdatePassword(model *pb.RegisterModel) (*pb.SellerInfo, error) {
	timeNow := time.Now()
	model.Password = encryPassword(model.Password)
	query := "update seller set password=$1,update_at=$2 where mobile=$3 returning id,nickname"
	seller := &pb.SellerInfo{}
	err := DB.QueryRow(query, model.Password, timeNow, model.Mobile).Scan(&seller.Id, &seller.Username)
	if err != nil {
		return nil, err
	}
	return seller, nil
}

//private functions
func encryPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password)) // 需要加密的字符串为 sharejs.com
	result := hex.EncodeToString(h.Sum(nil))
	return result
}
