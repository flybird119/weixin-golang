package db

import (
	"database/sql"
	"fmt"
	"strings"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//增加用户地址
func AddAddress(address *pb.AddressInfo) error {
	//首先获取用户的地址集合
	hasEquelsAddr := false
	isDefault := true
	addresses, err := FindAddressByUser(address)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	for i := 0; i < len(addresses); i++ {
		findAddress := addresses[i]
		//查看是否有相同的地址
		if address.Name == findAddress.Name && address.Tel == findAddress.Tel && address.Address == findAddress.Address {
			hasEquelsAddr = true
		}
		isDefault = false
	}
	//如果有相同的地址,那么返回
	if hasEquelsAddr {
		return nil
	}
	query := "insert into address (name,tel,address,user_id,is_default) values($1,$2,$3,$4,$5) returning id"
	log.Debugf("insert into address (name,tel,address,user_id,is_default) values('%s','%s','%s','%s',%s)", address.Name, address.Tel, address.Address, address.UserId, isDefault)
	err = DB.QueryRow(query, address.Name, address.Tel, address.Address, address.UserId, isDefault).Scan(&address.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//更新用户地址
func UpdateAddress(address *pb.AddressInfo) error {
	query := "update address set update_at=now() "

	var args []interface{}
	var condition string

	if address.Name != "" {
		args = append(args, address.Name)
		condition += fmt.Sprintf(",name=$%d", len(args))
	}

	if address.Tel != "" {
		args = append(args, address.Tel)
		condition += fmt.Sprintf(",tel=$%d", len(args))
	}

	if address.Address != "" {
		args = append(args, address.Address)
		condition += fmt.Sprintf(",address=$%d", len(args))
	}

	if address.SetDefault == 1 {
		query1 := fmt.Sprintf("update address set is_default=false where user_id='%s'", address.UserId)
		log.Debug(query1)
		_, err := DB.Exec(query1)
		if err != nil {
			misc.LogErr(err)
			return err
		}
		args = append(args, true)
		condition += fmt.Sprintf(",is_default=$%d", len(args))
	}
	args = append(args, address.UserId)
	condition += fmt.Sprintf(" where user_id=$%d", len(args))
	args = append(args, address.Id)
	condition += fmt.Sprintf(" and id=$%d", len(args))

	query += condition
	log.Debugf(query+" args:%#v", args)
	_, err := DB.Exec(query, args...)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//删除用户地址
func DelAddress(addresses []*pb.AddressInfo, userId string) error {
	query := fmt.Sprintf("delete from address where id in (${ids}) and user_id='%s'", userId)
	var idArray []interface{}
	if len(addresses) > 0 {
		query = strings.Replace(query, "${"+"ids"+"}",
			strings.Repeat(",'%s'", len(addresses))[1:], -1)
		for _, s := range addresses {
			idArray = append(idArray, s.Id)
		}
		query = fmt.Sprintf(query, idArray...)
		log.Debug(query)
	}

	_, err := DB.Exec(query)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//获取用户的地址
func FindAddressByUser(address *pb.AddressInfo) (addresses []*pb.AddressInfo, err error) {
	query := "select id,name,tel,address,user_id,is_default from address where user_id=$1"
	log.Debugf("select id,name,tel,address,user_id,is_default from address where user_id='%s'", address.UserId)
	rows, err := DB.Query(query, address.UserId)
	if err == sql.ErrNoRows {

		return addresses, nil
	}
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		searchAddress := &pb.AddressInfo{}
		addresses = append(addresses, searchAddress)
		err = rows.Scan(&searchAddress.Id, &searchAddress.Name, &searchAddress.Tel, &searchAddress.Address, &searchAddress.UserId, &searchAddress.IsDefault)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
	}
	return
}
