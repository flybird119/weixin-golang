package db

import (
	"database/sql"
	"errors"
	"fmt"

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
func HasExistAcoount(item *pb.AccountItem) (bool, error) {
	item.Id = ""
	query := "select id from account_item where store_id=$1 and order_id=$2 and item_type=$3"
	log.Debugf("select id from account_item where store_id='%s' and order_id='%s' and item_type=%d", item.StoreId, item.OrderId, item.ItemType)
	err := DB.QueryRow(query, item.StoreId, item.OrderId, item.ItemType).Scan(&item.Id)

	if err == sql.ErrNoRows || item.Id == "" {
		return false, nil
	} else if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return false, err
	}
	return true, nil
}

//store_id unsettled_balance
func ChangeAccountWithdrawalFee(account *pb.Account) error {
	query := "update account set unsettled_balance=unsettled_balance+$1,update_at=now() where store_id=$2 returning unsettled_balance ,id"
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, account.UnsettledBalance, account.StoreId).Scan(&account.UnsettledBalance, &account.Id)
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	tx.Commit()
	return nil
}

//修改账户可提现余额
func ChangAccountBalance(account *pb.Account) error {
	query := "update account set balance=balance+$1,update_at=now() where store_id=$2 returning balance,id"
	//开启事务
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, account.Balance, account.StoreId).Scan(&account.Balance, &account.Id)
	if err != nil {
		log.Error(err)
		return err
	}
	tx.Commit()
	return nil
}

//修改账户可提现余额
func ChangAccountBalanceWithTx(tx *sql.Tx, account *pb.Account) error {
	query := "update account set balance=balance+$1,update_at=now() where store_id=$2 returning balance,id"
	//开启事务
	log.Debug(query+" args:%s,%d", account.StoreId, account.Balance)
	err := tx.QueryRow(query, account.Balance, account.StoreId).Scan(&account.Balance, &account.Id)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug("hello")
	if account.Balance < 0 {
		return errors.New("sellerNoMoney")
	}
	return nil
}

//增加什么AccountItem
func AddAccountItemWithTx(tx *sql.Tx, item *pb.AccountItem) error {
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
	query = "insert into account_item (user_type,store_id,order_id,remark,item_type,item_fee,account_balance) values ($1,$2,$3,$4,$5,$6,$7)"
	log.Debugf("insert into account_item (user_type,store_id,order_id,remark,item_type,item_fee,account_balance) values (%d,'%s','%s','%s',%d,%d,%d)", item.UserType, item.StoreId, item.OrderId, item.Remark, item.ItemType, item.ItemFee, item.AccountBalance)
	_, err = tx.Exec(query, item.UserType, item.StoreId, item.OrderId, item.Remark, item.ItemType, item.ItemFee, item.AccountBalance)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//查找账户项列表
func FindAccountItems(submitModel *pb.FindAccountitemReq) (respModel pb.FindAccountitemResp, err error) {

	//>1	商家 待结算-交易完成
	//>2 商家 待结算-手续费
	//>4 商家 待结算-交易收入
	//>17 商家 可提现-交易完成
	//>18 商家 可提现-充值
	//>20 商家 可提现-体现
	//>24 商家 可提现-售后
	//>80 商家 所有待结算分类
	//>81 商家 所有可提现分类
	query := "select store_id,order_id,remark,item_type,item_fee,account_balance,extract(epoch from create_at)::bigint from account_item where 1=1"
	queryCount := "select count(*) from account_item where 1=1"
	querySumPlus := "select sum(item_fee) from account_item where item_fee>0"
	querySumReduce := "select sum(item_fee) from account_item where item_fee<0"

	var condition string
	if submitModel.StartAt != 0 && submitModel.EndAt != 0 {
		condition += fmt.Sprintf(" and extract(epoch from create_at)::bigint between %d and %d", submitModel.StartAt, submitModel.EndAt)
	}
	if submitModel.Type == 80 {
		condition += fmt.Sprintf(" and item_type in (1,2,4)")
	} else if submitModel.Type == 81 {
		condition += fmt.Sprintf(" and item_type in (17,18,20,24)")
	} else {
		condition += fmt.Sprintf(" and item_type=%d", submitModel.Type)
	}

	condition += fmt.Sprintf(" and store_id='%s'", submitModel.StoreId)
	//计算总个数
	queryCount += condition
	log.Debug(queryCount)
	err = DB.QueryRow(queryCount).Scan(&respModel.TotalCount)
	if err != nil {
		log.Error(err)
		return
	}
	if respModel.TotalCount == 0 {
		return
	}
	var sumplus, sumreduce sql.NullInt64
	//计算当前条件总收入
	querySumPlus += condition
	log.Debug(querySumPlus)
	err = DB.QueryRow(querySumPlus).Scan(&sumplus)
	if err != nil {
		log.Error(err)
		return
	}
	//计算当前条件总支出
	querySumReduce += condition
	log.Debug(querySumReduce)
	err = DB.QueryRow(querySumReduce).Scan(&sumreduce)
	if err != nil {
		log.Error(err)
		return
	}
	if sumplus.Valid {
		respModel.TotalIncome = sumplus.Int64
	}
	if sumreduce.Valid {
		respModel.TotalExpense = sumreduce.Int64
	}
	//遍历

	condition += " order by create_at desc"
	if submitModel.Page <= 0 {
		submitModel.Page = 1
	}
	if submitModel.Size <= 0 {
		submitModel.Size = 10
	}
	condition += fmt.Sprintf(" OFFSET %d LIMIT %d ", (submitModel.Page-1)*submitModel.Size, submitModel.Size)
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &pb.AccountItem{}
		respModel.Data = append(respModel.Data, item)
		//store_id,order_id,remark,item_type,item_fee,account_balance,extract(epoch from create_at)::bigint from account_item
		err = rows.Scan(&item.StoreId, &item.OrderId, &item.Remark, &item.ItemType, &item.ItemFee, &item.AccountBalance, &item.CreateAt)
		if err != nil {
			log.Error(err)
			return
		}
	}

	return
}

//获取账户详情
func GetAccountDetail(account *pb.Account) error {
	query := "select balance,unsettled_balance from account where store_id='%s'"
	query = fmt.Sprintf(query, account.StoreId)
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&account.Balance, &account.UnsettledBalance)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
