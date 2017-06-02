package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//根据手机号和密码获取用户
func GetUserByPasswordAndMobile(master *pb.Master) error {

	query := "select id from master where mobile='%s' and password='%s'"
	master.Password = encryPassword(master.Password)
	query = fmt.Sprintf(query, master.Mobile, master.Password)
	log.Debugf(query)
	err := DB.QueryRow(query).Scan(&master.Id)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

//提现列表
func WithdrawList(model *pb.StoreWithdrawalsModel) (models []*pb.StoreWithdrawalsModel, err error, totalCount int64) {
	query := "select w.id,w.store_id,w.card_type,w.card_no,w.card_name,w.username,w.withdraw_fee,w.status,w.apply_phone,extract(epoch from w.apply_at)::bigint,extract(epoch from w.complete_at)::bigint ,extract(epoch from w.accept_at)::bigint from withdrawals w where 1=1"
	queryCount := "select count(*) from withdrawals w where 1=1"
	var condition string
	if model.Status != 0 {
		condition += fmt.Sprintf(" and w.status=%d", model.Status)
	}

	if model.Id != "" {
		condition += fmt.Sprintf(" and w.id='%s'", model.Id)
	}

	if model.Page <= 0 {
		model.Page = 1
	}
	if model.Size <= 0 {
		model.Size = 10
	}
	//检查条数
	queryCount += condition
	log.Debugf(queryCount)
	err = DB.QueryRow(queryCount).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}

	if totalCount == 0 {
		return
	}
	//检索数据
	condition += fmt.Sprintf(" order by update_at desc OFFSET %d LIMIT %d", (model.Page-1)*model.Size, model.Size)
	query += condition
	log.Debugf(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		log.Error(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		withdraw := &pb.StoreWithdrawalsModel{}
		models = append(models, withdraw)
		var complete_at, accept_at sql.NullInt64
		// w.id,w.card_type,w.card_no,w.card_name,w.username,w.withdraw_fee,w.status,w.apply_phone,extract(epoch from w.apply_at)::bigint,extract(epoch from w.complete_at)::bigint extract(epoch from w.accept_at)
		err = rows.Scan(&withdraw.Id, &withdraw.StoreId, &withdraw.CardType, &withdraw.CardNo, &withdraw.CardName, &withdraw.Username, &withdraw.WithdrawFee, &withdraw.Status, &withdraw.ApplyPhone, &withdraw.ApplyAt, &complete_at, &accept_at)
		if err != nil {
			log.Error(err)
			return
		}
		if complete_at.Valid {
			withdraw.CompleteAt = complete_at.Int64
		}

		if accept_at.Valid {
			withdraw.AcceptAt = accept_at.Int64
		}

		query = fmt.Sprintf("select s.name,a.balance from store s join account a on s.id=a.store_id where s.id='%s'", withdraw.StoreId)
		log.Debug(query)
		err = DB.QueryRow(query).Scan(&withdraw.StoreName, &withdraw.Balance)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

//根据id获取提现信息
func GetWithdrawById(withdrawId string) (*pb.StoreWithdrawalsModel, error) {

	query := "select w.id,w.card_type,w.card_no,w.card_name,w.username,w.withdraw_fee,w.status,w.apply_phone,extract(epoch from w.apply_at)::bigint,extract(epoch from w.complete_at)::bigint ,extract(epoch from w.accept_at)::bigint from withdrawals w where w.id='%s' "
	query = fmt.Sprintf(query, withdrawId)
	log.Debug(query)
	withdraw := &pb.StoreWithdrawalsModel{}
	var complete_at, accept_at sql.NullInt64
	err := DB.QueryRow(query).Scan(&withdraw.Id, &withdraw.CardType, &withdraw.CardNo, &withdraw.CardName, &withdraw.Username, &withdraw.WithdrawFee, &withdraw.Status, &withdraw.ApplyPhone, &withdraw.ApplyAt, &complete_at, &accept_at)

	if complete_at.Valid {
		withdraw.CompleteAt = complete_at.Int64
	}

	if accept_at.Valid {
		withdraw.AcceptAt = accept_at.Int64
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return withdraw, nil
}

//开始处理提现
func UpdateWithdraw(model *pb.StoreWithdrawalsModel) error {
	query := "update withdrawals set update_at=now() "
	var condition string
	if model.Status != 0 {
		condition += fmt.Sprintf(",status=%d", model.Status)
	}
	if model.AcceptAt != 0 {
		condition += ",accept_at=now()"
	}
	if model.CompleteAt != 0 {
		condition += ",complete_at=now()"
	}
	condition += fmt.Sprintf(" where id='%s'", model.Id)
	query += condition
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

//private functions
func encryPassword(password string) string {
	h := md5.New()
	h.Write([]byte(password)) // 需要加密的字符串为 sharejs.com
	result := hex.EncodeToString(h.Sum(nil))
	return result
}
