package db

import (
	"database/sql"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//InitAccount 初始化account
func InitAccount(account *pb.Account) error {
	query := "select id from account where type=$1 and store_id=$2"
	log.Debugf("select id from account where type=%d and store_id=%s", account.Type, account.StoreId)
	err := DB.QueryRow(query, account.Type, account.StoreId).Scan(&account.Id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if account.Id == "" {
		query = "insert into account (type,store_id) values($1,$2)"
		log.Debugf("insert into account (type,store_id) values(%d,'%s')", account.Type, account.StoreId)
		_, err = DB.Exec(query, account.Type, account.StoreId)
		if err != nil {
			misc.LogErr(err)
			return err
		}
	}
	return nil
}

//增加什么AccountItem
func AddAccountItem(item *pb.AccountItem) error {
	query := "select id from account_item where user_type=$1 and store_id=$2 and order_id=$3 and item_type=$4"
	log.Debugf("select id from account_item where user_type=%d and store_id='%s' and order_id='%s' and item_type=%d", item.UserType, item.StoreId, item.OrderId, item.ItemType)
	err := DB.QueryRow(query, item.UserType, item.StoreId, item.OrderId, item.ItemType).Scan(&item.Id)
	if err != nil && err != sql.ErrNoRows {
		misc.LogErr(err)
		return err
	}
	if item.Id != "" {

		return nil
	}
	tx, err := DB.Begin()
	if err != nil {
		misc.LogErr(err)
		return err
	}
	defer tx.Rollback()
	query = "insert into account_item (user_type,store_id,order_id,remark,item_type,item_fee,account_balance) values ($1,$2,$3,$4,$5,$6,$7)"
	log.Debugf("insert into account_item (user_type,store_id,order_id,remark,item_type,item_fee,account_balance) values (%d,'%s','%s','%s',%d,%d,%d)", item.UserType, item.StoreId, item.OrderId, item.Remark, item.ItemType, item.ItemFee, item.AccountBalance)
	_, err = tx.Exec(query, item.UserType, item.StoreId, item.OrderId, item.Remark, item.ItemType, item.ItemFee, item.AccountBalance)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	tx.Commit()
	return nil
}

//检查账户是否存在
func hasExistAcoount(item *pb.AccountItem) {

}

//store_id unsettled_balance
func ChangeAccountWithdrawalFee(account *pb.Account) error {
	query := "update account set unsettled_balance=unsettled_balance+$1,update_at=now() where store_id=$2 returning unsettled_balance ,id"
	tx, err := DB.Begin()
	if err != nil {
		misc.LogErr(err)
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, account.UnsettledBalance, account.StoreId).Scan(&account.UnsettledBalance, &account.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	tx.Commit()
	return nil

}
