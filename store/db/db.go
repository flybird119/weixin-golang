package db

import (
	"database/sql"
	"time"

	"github.com/goushuyun/weixin-golang/seller/role"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//AddStore 通过手机号和登录密码检查商家是否存在
func AddStore(store *pb.Store) error {
	query := "insert into store (name,status,expire_at) values($1,$2,$3) returning id,extract(epoch from create_at)::integer "
	now := time.Now()
	now = now.Add(7 * 24 * time.Hour)
	err := DB.QueryRow(query, store.Name, pb.StoreStatus_Normal, now).Scan(&store.Id, &store.CreateAt)
	store.ExpireAt = now.Unix()
	if err != nil {
		return err
	}
	return nil
}

//UpdateStore 增加商家和店铺的映射
func UpdateStore(store *pb.Store) error {
	query := "update store set name=$1,profile=$2 where id=$3"
	_, err := DB.Query(query, store.Name, store.Profile, store.Id)
	log.Debugf("==============>update store set name=%s,profile=%s where id=%s", store.Name, store.Profile, store.Id)
	if err != nil {
		return err
	}
	return nil
}

//AddStoreSellerMap 增加商家和店铺的映射
func AddStoreSellerMap(store *pb.Store, role int64) error {
	query := "insert into map_store_seller (seller_id,store_id,role) values($1,$2,$3)"
	_, err := DB.Query(query, store.Seller.Id, store.Id, role)
	if err != nil {
		return err
	}
	return nil
}

//AddRealStore 增加实体店
func AddRealStore(realStore *pb.RealStore) error {
	query := "insert into real_shop (name,province_code,city_code,scope_code,address,images,store_id) values($1,$2,$3,$4,$5,$6,$7) returning id"
	err := DB.QueryRow(query, realStore.Name, realStore.ProvinceCode, realStore.CityCode, realStore.ScopeCode, realStore.Address, realStore.Images, realStore.StoreId).Scan(&realStore.Id)
	if err != nil {
		return err
	}
	return nil
}

//UpdateRealStore 修改实体店铺信息
func UpdateRealStore(realStore *pb.RealStore) error {
	query := "update real_shop set name=$1,province_code=$2,city_code=$3,scope_code=$4,address=$5,images=$6 where id=$7"
	log.Debugf("update real_shop set name=%s,province_code=%s,city_code=%s,scope_code=%s,address=%s,images=%s where id=%s", realStore.Name, realStore.ProvinceCode, realStore.CityCode, realStore.ScopeCode, realStore.Address, realStore.Images, realStore.Id)
	_, err := DB.Query(query, realStore.Name, realStore.ProvinceCode, realStore.CityCode, realStore.ScopeCode, realStore.Address, realStore.Images, realStore.Id)
	if err != nil {
		return err
	}
	return nil
}

//GetStoreInfo 获取店铺的信息
func GetStoreInfo(store *pb.Store) error {
	//获取店铺基本信息
	query := "select name,logo,status,profile,extract(epoch from expire_at)::integer,address,business_license,extract(epoch from create_at)::integer from store where id=$1"
	var logo, profile, address, businessLicense sql.NullString
	err := DB.QueryRow(query, store.Id).Scan(&store.Name, &logo, &store.Status, &profile, &store.ExpireAt, &address, &businessLicense, &store.CreateAt)

	log.Debugf("select name,logo,status,profile,extract(epoch from expire_at)::integer,address,business_license,extract(epoch from create_at)::integer from store where id='%s'", store.Id)

	if err != nil {
		log.Debugf("Err:%s !!select name,logo,status,profile,service_mobiles,extract(epoch from s.expire_at)::integer,address,business_license,extract(epoch from s.create_at)::integer where id=%s", err, store.Id)
		return err
	}
	if logo.Valid {
		store.Logo = logo.String
	}
	if profile.Valid {
		store.Profile = profile.String
	}
	if address.Valid {
		store.Address = address.String
	}
	if businessLicense.Valid {
		store.BusinessLicense = businessLicense.String
	}

	//获取店铺负责人手机号
	query = "select s.mobile from seller s join map_store_seller ms on ms.seller_id=s.id where ms.role=$1 and ms.store_id=$2"
	log.Debugf("select s.mobile from seller s join map_store_seller ms on ms.seller_id=s.id where ms.role=%s and ms.store_id=%s", role.InterAdmin, store.Id)
	err = DB.QueryRow(query, role.InterAdmin, store.Id).Scan(&store.AdminMobile)
	if err != nil {
		log.Errorf("select s.mobile from seller s join map_store_seller ms on ms.seller_id=s.id where ms.role=%d and ms.store_id=%s", role.InterAdmin, store.Id)
	}
	return nil
}

//GetSellerStoreRole 获取商户权限
func GetSellerStoreRole(sellerId, storeId string) (int64, error) {
	query := "select role from map_store_seller where seller_id=$1 and store_id=$2"
	log.Debugf("select role from map_store_seller where seller_id=%s and store_id=%s", sellerId, storeId)
	var role int64
	err := DB.QueryRow(query, sellerId, storeId).Scan(&role)
	if err != nil {
		log.Errorf("select role from map_store_seller where seller_id=%s and store_id=%s", sellerId, storeId)
		return 0, err
	}
	return role, nil
}

//ChangeStoreLogo 修改店铺头像
func ChangeStoreLogo(image, store_id string) error {
	query := "update store set logo=$1 where id=$2"
	log.Debugf("update store set logo=%s where id=%s", image, store_id)
	_, err := DB.Exec(query, image, store_id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//GetStoreShops 获取店铺下的实习店列表
func GetStoreShops(storeId string) (r []*pb.RealStore, err error) {
	query := "select id,name,province_code,city_code,scope_code,address,images ,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer from real_shop where store_id=$1"
	log.Debugf("select id ,name,province_code,city_code,scope_code,address,images ,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer from real_shop where store_id=%s", storeId)
	rows, err := DB.Query(query, storeId)
	if err != nil {
		log.Errorf("select id ,name,province_code,city_code,scope_code,address,images ,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer from real_shop where store_id=%s", storeId)
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var store pb.RealStore
		r = append(r, &store)
		err = rows.Scan(&store.Id, &store.Name, &store.ProvinceCode, &store.CityCode, &store.ScopeCode, &store.Address, &store.Images, &store.CreateAt, &store.UpdateAt)

		if err != nil {
			log.Debug(err)
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		log.Debug("scan rows err last error: %s", err)
		return nil, err
	}
	return r, nil
}

//TransferStore 转让店铺
func TransferStore(sellerId, storeId string) error {
	query := "update map_store_seller set seller_id=$1 where store_id=$2 and role=$3"
	log.Debugf("update map_store_seller set seller_id=%s where store_id=%s and role=%d", sellerId, storeId, role.InterAdmin)
	_, err := DB.Exec(query, sellerId, storeId, role.InterAdmin)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
