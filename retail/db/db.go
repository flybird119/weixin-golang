package db

import (
	"database/sql"
	"errors"

	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	sellerDB "github.com/goushuyun/weixin-golang/seller/db"
	"github.com/wothing/log"
)

//增加零售
func AddRetail(retail *pb.RetailSubmitModel) error {
	//首先先插入retail
	query := "insert into retail (total_fee,store_id,school_id,handle_staff_id,goods_fee) values($1,$2,$3,$4,0) returning id"
	log.Debugf("insert into retail (total_fee,store_id,school_id,handle_staff_id,goods_fee) values(%d,'%s','%s','%s',0) returning id", retail.TotalFee, retail.StoreId, retail.SchoolId, retail.SellerId)
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRow(query, retail.TotalFee, retail.StoreId, retail.SchoolId, retail.SellerId).Scan(&retail.RetailId)
	if err != nil {
		log.Error(err)
		return err
	}
	//增加零售项
	hasErr := false
	for i := 0; i < len(retail.Items); i++ {
		retail.Items[i].RetailId = retail.RetailId
		retail.Items[i].HasStock = true
		err := AddRetailItem(tx, retail.Items[i])
		if err != nil {
			log.Error(err)
			hasErr = true
		}
	}
	if hasErr {
		return errors.New("noStock")
	}
	tx.Commit()
	return nil
}

//增加零售项
func AddRetailItem(tx *sql.Tx, item *pb.RetailItem) (err error) {
	//插入零售项
	var query string
	//修改商品数量
	if item.Type == 0 {
		query = "update goods set new_book_amount=new_book_amount-$1 where id=$2 returning new_book_amount,new_book_price"
	} else {
		query = "update goods set old_book_amount=old_book_amount-$1 where id=$2 returning old_book_amount,old_book_price"
	}
	log.Debugf(query+" amount:%d,goodsId=%s", item.Amount, item.GoodsId)

	var amount, price int64
	err = tx.QueryRow(query, item.Amount, item.GoodsId).Scan(&amount, &price)
	if err != nil {
		log.Error(err)
		return
	}
	//检测数量
	item.CurrentAmount = amount + item.Amount
	if amount < 0 {
		item.HasStock = false
		err = errors.New("noStock")
		return
	}

	query = "insert into retail_item (goods_id,retail_id,type,amount,price) values($1,$2,$3,$4,$5)"
	log.Debugf("insert into retail_item (goods_id,retail_id,type,amount,price) values('%s','%s','%d','%d','%d')", item.GoodsId, item.RetailId, item.Type, item.Amount, price)
	_, err = tx.Exec(query, item.GoodsId, item.RetailId, item.Type, item.Amount, price)
	if err != nil {
		log.Error(err)
		return
	}
	//修改零售商品费用
	itemFee := item.Amount * price
	query = "update retail set goods_fee=goods_fee+$1 where id=$2"
	log.Debugf("update retail set goods_fee=goods_fee+%d where id='%s'", itemFee, item.RetailId)
	_, err = tx.Exec(query, itemFee, item.RetailId)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

//售后检索
func FindRetails(retail *pb.Retail) (details []*pb.RetailDetail, err error, totalCount int64) {
	query := "select r.id,r.total_fee,r.store_id,r.school_id,r.handle_staff_id,extract(epoch from r.create_at)::integer,r.goods_fee from retail r where 1=1"
	countQuery := "select count(*) from retail r where 1=1"
	var args []interface{}
	var condition string
	if retail.SchoolId != "" {
		args = append(args, retail.SchoolId)
		condition += fmt.Sprintf(" and r.school_id=$%d", len(args))
	}

	if retail.Isbn != "" {
		args = append(args, retail.Isbn)
		condition += fmt.Sprintf(" and (exists (select * from retail_item ri join  goods g on ri.goods_id=g.id where ri.retail_id=r.id and g.isbn=$%d))", len(args))

	}

	if retail.StartAt != 0 && retail.EndAt != 0 {
		args = append(args, retail.StartAt)
		condition += fmt.Sprintf(" and extract(epoch from r.create_at)::integer between $%d and $%d", len(args), len(args)+1)
		args = append(args, retail.EndAt)
	}

	args = append(args, retail.StoreId)
	condition += fmt.Sprintf(" and r.store_id=$%d", len(args))

	countQuery += condition
	if retail.Page <= 0 {
		retail.Page = 1
	}
	if retail.Size <= 0 {
		retail.Size = 10
	}

	condition += fmt.Sprintf("  order by update_at desc OFFSET %d LIMIT %d ", (retail.Page-1)*retail.Size, retail.Size)

	query += condition
	log.Debugf(query+" args:%#v", args)

	//统计满足条件的总条数
	log.Debugf(countQuery+" args:%#v", args)
	err = DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}
	if totalCount <= 0 {
		return
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		detail := &pb.RetailDetail{}
		findRetail := &pb.Retail{}
		detail.Retail = findRetail
		// r.id,r.total_fee,r.store_id,r.school_id,r.handle_staff_id,extract(epoch from r.create_at)::integer
		err = rows.Scan(&findRetail.Id, &findRetail.TotalFee, &findRetail.StoreId, &findRetail.SchoolId, &findRetail.HandleStaffId, &findRetail.CreateAt, &findRetail.GoodsFee)
		if err != nil {
			log.Error(err)
			return
		}
		//获取items
		items, err := GetRetailItems(findRetail)
		if err != nil {
			log.Error(err)
			return details, err, totalCount
		}
		detail.Items = items
		//获取处理人的信息
		sellerInfo, err := sellerDB.GetSellerById(findRetail.HandleStaffId)
		if err != nil {
			log.Error(err)
			return details, err, totalCount
		}
		detail.ChargeMan = sellerInfo

		details = append(details, detail)

	}

	return
}

//
func GetRetailItems(retail *pb.Retail) (items []*pb.RetailItem, err error) {
	query := "select ri.id,g.id,ri.type,ri.amount,ri.price,b.title,b.isbn,b.image,b.price from retail_item ri join goods g on ri.goods_id=g.id join books b on g.book_id=b.id where retail_id='%s'"
	query = fmt.Sprintf(query, retail.Id)
	log.Debugf(query)
	rows, err := DB.Query(query)
	//如果出现无结果异常
	if err == sql.ErrNoRows {
		return items, nil
	}
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	//遍历搜索结果
	for rows.Next() {
		item := &pb.RetailItem{}
		items = append(items, item)
		err = rows.Scan(&item.Id, &item.GoodsId, &item.Type, &item.Amount, &item.Price, &item.BookTitle, &item.BookIsbn, &item.BookImage, &item.OriginPrice)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
	}
	return
}
