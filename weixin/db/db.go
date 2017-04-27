package db

import (
	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
)

func SaveAccount(accout *pb.GetAuthBaseInfoResp) error {
	query := "insert into official_accounts(nick_name, head_img, user_name, principal_name, qrcode_url, service_type_info, verify_type_info, appid) values($1, $2, $3, $4, $5, $6, $7, $8)"

	log.Debugf("insert into official_accounts(nick_name, head_img, user_name, principal_name, qrcode_url, service_type_info, verify_type_info, appid) values('%s', '%s', '%s', '%s', '%s', %d, %d, '%s')", accout.AuthorizerInfo.NickName, accout.AuthorizerInfo.HeadImg, accout.AuthorizerInfo.UserName, accout.AuthorizerInfo.PrincipalName, accout.AuthorizerInfo.QrcodeUrl, accout.AuthorizerInfo.ServiceTypeInfo, accout.AuthorizerInfo.VerifyTypeInfo, accout.AuthorizationInfo.AuthorizerAppid)

	_, err := DB.Exec(query, accout.AuthorizerInfo.NickName, accout.AuthorizerInfo.HeadImg, accout.AuthorizerInfo.UserName, accout.AuthorizerInfo.PrincipalName, accout.AuthorizerInfo.QrcodeUrl, accout.AuthorizerInfo.ServiceTypeInfo, accout.AuthorizerInfo.VerifyTypeInfo, accout.AuthorizationInfo.AuthorizerAppid)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func SaveAppidToStore(store_id, app_id string) error {
	query := "update store set appid = $1 where id = $2"
	log.Debugf("update store set appid = '%s' where id = '%s'", app_id, store_id)

	_, err := DB.Exec(query, app_id, store_id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
