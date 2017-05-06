package db

import (
	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
)

func GetAccountInfoByStoreId(store_id string) (*pb.OfficialAccount, error) {
	account := &pb.OfficialAccount{}

	query := "select oa.id, oa.appid, oa.nick_name, oa.head_img, oa.user_name, oa.principal_name, oa.qrcode_url, oa.service_type_info, oa.verify_type_info, extract(epoch from store.create_at)::integer create_at, store.authorizer_refresh_token, oa.wechat_id from official_accounts as oa, store where store.id = $1 and store.appid = oa.appid"

	log.Debugf("select oa.id, oa.appid, oa.nick_name, oa.head_img, oa.user_name, oa.principal_name, oa.qrcode_url, oa.service_type_info, oa.verify_type_info, extract(epoch from store.create_at)::integer create_at, store.authorizer_refresh_token, oa.wechat_id from official_accounts as oa,  store where store.id = '%s' and store.appid = oa.appid", store_id)

	err := DB.QueryRow(query, store_id).Scan(&account.Id, &account.Appid, &account.NickName, &account.HeadImg, &account.UserName, &account.PrincipalName, &account.QrcodeUrl, &account.ServiceTypeInfo, &account.VerifyTypeInfo, &account.CreateAt, &account.RefreshToken, &account.WechatId)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return account, nil
}

func SaveAccount(accout *pb.GetAuthBaseInfoResp) error {
	query := "insert into official_accounts(nick_name, head_img, user_name, principal_name, qrcode_url, service_type_info, verify_type_info, appid, wechat_id) values($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	log.Debugf("insert into official_accounts(nick_name, head_img, user_name, principal_name, qrcode_url, service_type_info, verify_type_info, appid, wechat_id) values('%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s')", accout.AuthorizerInfo.NickName, accout.AuthorizerInfo.HeadImg, accout.AuthorizerInfo.UserName, accout.AuthorizerInfo.PrincipalName, accout.AuthorizerInfo.QrcodeUrl, accout.AuthorizerInfo.ServiceTypeInfo.Id, accout.AuthorizerInfo.VerifyTypeInfo.Id, accout.AuthorizationInfo.AuthorizerAppid, accout.AuthorizerInfo.Alias)

	_, err := DB.Exec(query, accout.AuthorizerInfo.NickName, accout.AuthorizerInfo.HeadImg, accout.AuthorizerInfo.UserName, accout.AuthorizerInfo.PrincipalName, accout.AuthorizerInfo.QrcodeUrl, accout.AuthorizerInfo.ServiceTypeInfo.Id, accout.AuthorizerInfo.VerifyTypeInfo.Id, accout.AuthorizationInfo.AuthorizerAppid, accout.AuthorizerInfo.Alias)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func AccountExist(accout *pb.GetAuthBaseInfoResp) (bool, error) {
	query := "select count(*) from official_accounts where appid = $1"
	var total int64
	err := DB.QueryRow(query, accout.AuthorizationInfo.AuthorizerAppid).Scan(&total)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return total > 0, nil
}

func UpdateAccount(accout *pb.GetAuthBaseInfoResp) error {

	return nil
}

func SaveAuthorizerInfoToStore(store_id, app_id, token string) error {
	query := "update store set appid = $1, authorizer_refresh_token = $2 where id = $3"
	log.Debugf("update store set appid = '%s', authorizer_refresh_token = '%s' where id = '%s'", app_id, token, store_id)

	_, err := DB.Exec(query, app_id, token, store_id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
