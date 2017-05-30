package db

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//增加商品 book_id store_id isbn  goods.location
func AddGoods(goods *pb.Goods) error {
	var bookId, isbn string
	//首先根据isbn获取当前用户有没有保存goods
	query := "select id ,book_id,isbn from goods where isbn=$1 and store_id=$2"
	log.Debugf("select id ,book_id,isbn from goods where isbn=%s and store_id=$2", goods.Isbn, goods.StoreId)
	err := DB.QueryRow(query, goods.Isbn, goods.StoreId).Scan(&goods.Id, &bookId, &isbn)
	//如果检查失败
	if err == sql.ErrNoRows {
		//如果用户没有上传过改商品
		query = "insert into goods (book_id,store_id,isbn) values($1,$2,$3) returning id"
		log.Debugf("insert into goods (book_id,store_id,isbn) values(%s,%s,%s) returning id", goods.BookId, goods.StoreId, goods.Isbn)
		err = DB.QueryRow(query, goods.BookId, goods.StoreId, goods.Isbn).Scan(&goods.Id)
		if err != nil {
			log.Errorf("%+v", err)
			return err
		}

	} else if err != nil {
		log.Errorf("%+v", err)
		return err
	} else {
		//更改图书id
		if goods.BookId != bookId && goods.BookId != "" {
			err = UpdateGoods(&pb.Goods{Id: goods.Id, BookId: goods.BookId, Isbn: goods.Isbn})
			if err != nil {
				log.Errorf("%+v", err)
				return err
			}
		}
	}
	//遍历location
	for i := 0; i < len(goods.Location); i++ {
		goods.Location[i].GoodsId = goods.Id
		err = AddGoodsLocation(goods.Location[i])
		if err != nil {
			log.Errorf("%+v", err)
			return err
		}
	}
	return nil

}

//AddGoodsLoaction 增加货架位 goods_id  type storehouse_id shelf_id floor_id amount
func AddGoodsLocation(loc *pb.GoodsLocation) error {
	//首先查找货架位
	query := "select id from goods_location where goods_id=$1 and type=$2 and storehouse_id=$3 and shelf_id=$4 and floor_id=$5"
	log.Debugf("select id from goods_location where goods_id=%s and type=%d and storehouse_id=%s and shelf_id=%s and floor_id=%s", loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId)
	err := DB.QueryRow(query, loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId).Scan(&loc.Id)
	//如果检查失败
	if err == sql.ErrNoRows {
		//如果用户没有上传过货架位

		query = "insert into goods_location (goods_id,type,storehouse_id,shelf_id,floor_id) values($1,$2,$3,$4,$5)"

		log.Debugf("insert into goods_location (goods_id,type,storehouse_id,shelf_id,floor_id) values(%s,%d,%s,%s,%s)", loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId)
		_, err = DB.Exec(query, loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId)
		if err != nil {
			log.Errorf("%+v", err)
			return err
		}
	} else if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	//增加书本数量
	query = "update goods "
	debugQuery := "update goods"

	if loc.Type == 0 {
		query = query + " set new_book_amount=new_book_amount+$1,new_book_price=$2, has_new_book=true"
		debugQuery = debugQuery + " set new_book_amount=new_book_amount+%d,new_book_price=%d, has_new_book=true"
	} else if loc.Type == 1 {
		query = query + " set old_book_amount=old_book_amount+$1,old_book_price=$2, has_old_book=true"
		debugQuery = debugQuery + " set old_book_amount=old_book_amount+%d,old_book_price=%d,has_old_book=true"
	}
	updateTime := time.Now()
	//打开销售状态
	query = query + ",update_at=$3,is_selling=true where id=$4"
	debugQuery = debugQuery + ",update_at=%f,is_selling=true where id=%s"

	//修改时间
	log.Debugf(debugQuery, loc.Amount, loc.Price, updateTime, loc.GoodsId)
	_, err = DB.Exec(query, loc.Amount, loc.Price, updateTime, loc.GoodsId)
	if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	return nil
}

//UpdateGoods 更新商品 修改数量 book_id isbn title new_book_amount new_book_price old_book_price old_book_amount is_selling
func UpdateGoods(goods *pb.Goods) error {

	query := "update goods set update_at= now()"

	//动态拼接参数
	var args []interface{}
	var condition string
	if goods.BookId != "" {
		args = append(args, goods.BookId)
		condition += fmt.Sprintf(",book_id=$%d", len(args))
	}
	if goods.Isbn != "" {
		args = append(args, goods.Isbn)
		condition += fmt.Sprintf(",isbn=$%d", len(args))
	}
	if goods.NewBookAmount != 0 {
		if goods.NewBookAmount == -100 {
			args = append(args, 0)
			condition += fmt.Sprintf(",new_book_amount=$%d", len(args))
		} else {
			args = append(args, goods.NewBookAmount)
			condition += fmt.Sprintf(",new_book_amount=$%d", len(args))
		}

	}
	if goods.NewBookPrice != 0 {
		if goods.NewBookPrice == -100 {
			args = append(args, 0)
		} else {
			args = append(args, goods.NewBookPrice)
		}

		condition += fmt.Sprintf(",new_book_price=$%d", len(args))
	}
	if goods.OldBookAmount != 0 {
		if goods.OldBookAmount == -100 {
			args = append(args, 0)
			condition += fmt.Sprintf(",old_book_amount=$%d", len(args))
		} else {
			args = append(args, goods.OldBookAmount)
			condition += fmt.Sprintf(",old_book_amount=$%d", len(args))
		}

	}
	if goods.OldBookPrice != 0 {
		if goods.OldBookPrice == -100 {
			args = append(args, 0)
			condition += fmt.Sprintf(",old_book_price=$%d", len(args))
		} else {
			args = append(args, goods.OldBookPrice)
			condition += fmt.Sprintf(",old_book_price=$%d", len(args))
		}
	}
	if goods.SalesStatus != 0 {
		if goods.SalesStatus == -100 {
			args = append(args, true)
			condition += fmt.Sprintf(",is_selling=$%d", len(args))
		} else {
			args = append(args, false)
			condition += fmt.Sprintf(",is_selling=$%d", len(args))
		}
	}

	args = append(args, goods.Id)
	condition += fmt.Sprintf(" where id=$%d", len(args))
	log.Debugf(query+condition+"%+v", args)
	_, err := DB.Exec(query+condition, args...)
	if err != nil {
		log.Debugf("%+v", err)
		return err
	}

	return nil
}

//SearchGoods 搜索图书 isbn SearchAmount
func SearchGoods(goods *pb.Goods) (r []*pb.GoodsSearchResult, err error, totalCount int64) {
	query := "select %s from books b join goods g on b.id = g.book_id where 1=1 and is_selling=true "
	param := "b.id,b.store_id,b.title,b.isbn,b.price,b.author,b.publisher,b.pubdate,b.subtitle,b.image,b.summary,g.id, g.store_id,g.new_book_amount,g.new_book_price,g.old_book_amount,g.old_book_price,extract(epoch from g.create_at)::integer,extract(epoch from g.update_at)::integer,g.is_selling,g.has_new_book,g.has_old_book"
	query = fmt.Sprintf(query, param)
	queryCount := "select count(*) from books b join goods g on b.id = g.book_id where 1=1 and is_selling=true"
	//动态拼接参数
	var args []interface{}
	var condition string

	if goods.Isbn != "" {
		args = append(args, goods.Isbn)
		condition += fmt.Sprintf(" and b.isbn=$%d", len(args))
	}
	if goods.Author != "" {
		args = append(args, misc.FazzyQuery(goods.Author))
		condition += fmt.Sprintf(" and Lower(b.author) like Lower($%d)", len(args))
	}
	if goods.Publisher != "" {
		args = append(args, misc.FazzyQuery(goods.Publisher))
		condition += fmt.Sprintf(" and Lower(b.publisher) like Lower($%d)", len(args))
	}

	if goods.SearchType != -100 {
		if goods.SearchType == 0 {
			if goods.SearchAmount != 0 {
				if goods.SearchAmount == 1 {
					condition += " and has_new_book=true and g.new_book_amount<=0"
				} else {
					condition += " and has_new_book=true and g.new_book_amount>0"
				}
			} else {
				condition += " and has_new_book=true"
			}

		} else {
			if goods.SearchAmount != 0 {
				if goods.SearchAmount == 1 {
					condition += "  and has_old_book=true and g.old_book_amount<=0"
				} else {
					condition += " and has_old_book=true and g.old_book_amount>0"
				}
			} else {
				condition += " and has_old_book=true"
			}

		}
	} else {
		if goods.SearchAmount != 0 {
			if goods.SearchAmount == 1 {
				condition += "  and ((has_new_book=true and g.new_book_amount<=0) or (has_old_book=true and g.old_book_amount<=0))"
			} else {
				condition += "  and ((has_new_book=true and g.new_book_amount>0) or (has_old_book=true and g.old_book_amount>0))"
			}
		} else {
			condition += " and (has_old_book=true or has_new_book=true)"
		}

	}

	if goods.SearchPicture != -100 {
		if goods.SearchPicture == 0 {
			condition += " and b.image<>''"
		} else {
			condition += " and b.image=''"
		}
	}
	args = append(args, goods.StoreId)
	condition += fmt.Sprintf(" and g.store_id=$%d", len(args))

	if goods.Id != "" {
		args = append(args, goods.Id)
		condition += fmt.Sprintf(" and g.id=$%d", len(args))
	}

	if goods.Page <= 0 {
		goods.Page = 1
	}
	if goods.Size <= 0 {
		goods.Size = 20
	}

	if goods.Title != "" {
		args = append(args, misc.FazzyQuery(goods.Title))
		condition += fmt.Sprintf(" and Lower(b.title) like Lower($%d)", len(args))

		queryCount += condition
		err = DB.QueryRow(queryCount, args...).Scan(&totalCount)
		if err != nil {
			log.Error(err)
			return nil, err, totalCount
		}
		if totalCount == 0 {

			return
		}

		args = append(args, goods.Title)
		condition += fmt.Sprintf(" order by  Lower(title) <-> Lower($%d) ,g.update_at desc", len(args))
	} else {

		queryCount += condition
		err = DB.QueryRow(queryCount, args...).Scan(&totalCount)
		if err != nil {
			log.Error(err)
			return nil, err, totalCount
		}
		if totalCount == 0 {

			return
		}
		condition += " order by g.update_at desc"
		condition += fmt.Sprintf(" OFFSET %d LIMIT %d ", (goods.Page-1)*goods.Size, goods.Size)
	}
	query += condition
	log.Debugf(query+"%+v", args)
	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Error(err)
		return nil, err, totalCount
	}
	defer rows.Close()

	for rows.Next() {
		book := &pb.Book{}
		searchGoods := &pb.Goods{}

		var newbookModel *pb.GoodsSalesModel
		var oldbookModel *pb.GoodsSalesModel

		/**	param := "b.id,b.store_id,b.title,b.isbn,b.price,b.author,b.publisher,b.pubdate,b.subtitle,b.image,b.summary,g.id, g.store_id,g.new_book_amount,g.new_book_price,g.old_book_amount,g.old_book_price,extract(epoch from g.create_at)::integer,extract(epoch from g.update_at)::integer,g.is_selling"
		 */
		//遍历数据
		err = rows.Scan(&book.Id, &book.StoreId, &book.Title, &book.Isbn, &book.Price, &book.Author, &book.Publisher, &book.Pubdate, &book.Subtitle, &book.Image, &book.Summary, &searchGoods.Id, &searchGoods.StoreId, &searchGoods.NewBookAmount, &searchGoods.NewBookPrice, &searchGoods.OldBookAmount, &searchGoods.OldBookPrice, &searchGoods.CreateAt, &searchGoods.UpdateAt, &searchGoods.IsSelling, &searchGoods.HasNewBook, &searchGoods.HasOldBook)
		if err != nil {
			return nil, err, totalCount
		}

		if goods.SearchType == -100 {
			newLocations, _ := SearchGoodsLoaction(searchGoods.Id, 0)
			oldLocations, _ := SearchGoodsLoaction(searchGoods.Id, 1)

			if searchGoods.HasNewBook {
				newbookModel = &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 0, Price: searchGoods.NewBookPrice, Amount: searchGoods.NewBookAmount, Location: newLocations}
			}
			if searchGoods.HasOldBook {
				oldbookModel = &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 1, Price: searchGoods.OldBookPrice, Amount: searchGoods.OldBookAmount, Location: oldLocations}
			}

		} else {
			if goods.SearchType == 0 {
				newLocations, _ := SearchGoodsLoaction(searchGoods.Id, 0)
				if searchGoods.HasNewBook {
					newbookModel = &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 0, Price: searchGoods.NewBookPrice, Amount: searchGoods.NewBookAmount, Location: newLocations}

				}
			} else {
				oldLocations, _ := SearchGoodsLoaction(searchGoods.Id, 1)
				if searchGoods.HasOldBook {
					oldbookModel = &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 1, Price: searchGoods.OldBookPrice, Amount: searchGoods.OldBookAmount, Location: oldLocations}

				}
			}
		}

		r = append(r, &pb.GoodsSearchResult{Book: book, GoodsId: searchGoods.Id, StoreId: searchGoods.StoreId, UpdateAt: searchGoods.UpdateAt, NewBook: newbookModel, OldBook: oldbookModel})
	}

	return r, nil, totalCount
}

//SearchGoods 搜索图书 isbn 用于用户端搜索
func SearchGoodsNoLocation(goods *pb.Goods) (r []*pb.GoodsSearchResult, err error) {
	query := "select %s from books b join goods g on b.id = g.book_id where 1=1 and is_selling=true and (g.has_new_book=true or g.has_old_book=true)"
	param := "b.id,b.store_id,b.title,b.isbn,b.price,b.author,b.publisher,b.pubdate,b.subtitle,b.image,b.summary,g.id, g.store_id,g.new_book_amount,g.new_book_price,g.old_book_amount,g.old_book_price,extract(epoch from g.create_at)::integer,extract(epoch from g.update_at)::integer,g.is_selling,g.has_new_book,g.has_old_book"
	query = fmt.Sprintf(query, param)
	//动态拼接参数
	var args []interface{}
	var condition string

	if goods.Isbn != "" {
		args = append(args, goods.Isbn)
		condition += fmt.Sprintf(" and b.isbn=$%d", len(args))
	}

	args = append(args, goods.StoreId)
	condition += fmt.Sprintf(" and g.store_id=$%d", len(args))

	if goods.Id != "" {
		args = append(args, goods.Id)
		condition += fmt.Sprintf(" and g.id=$%d", len(args))
	}

	if goods.Page <= 0 {
		goods.Page = 1
	}
	if goods.Size <= 0 {
		goods.Size = 20
	}

	if goods.Title != "" {
		//查找有库存的图书
		condition += " and ((g.new_book_amount>0) or (g.old_book_amount>0))"
		args = append(args, misc.FazzyQuery(goods.Title))
		condition += fmt.Sprintf(" and ( Lower(b.title) like Lower($%d) or Lower(b.author) like Lower($%d) or Lower(b.publisher) like Lower($%d))", len(args), len(args)+1, len(args)+2)
		args = append(args, misc.FazzyQuery(goods.Title))
		args = append(args, misc.FazzyQuery(goods.Title))

		args = append(args, goods.Title)
		condition += fmt.Sprintf(" order by  Lower(b.title) <-> Lower($%d) ,g.update_at desc", len(args))
	} else {
		condition += " order by g.update_at desc"
		condition += fmt.Sprintf(" OFFSET %d LIMIT %d ", (goods.Page-1)*goods.Size, goods.Size)
	}
	query += condition
	log.Debugf(query+"%+v", args)
	rows, err := DB.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		book := &pb.Book{}
		searchGoods := &pb.Goods{}

		var newbookModel *pb.GoodsSalesModel
		var oldbookModel *pb.GoodsSalesModel
		var hasNewBook bool
		var hasOldBook bool
		//遍历数据
		err = rows.Scan(&book.Id, &book.StoreId, &book.Title, &book.Isbn, &book.Price, &book.Author, &book.Publisher, &book.Pubdate, &book.Subtitle, &book.Image, &book.Summary, &searchGoods.Id, &searchGoods.StoreId, &searchGoods.NewBookAmount, &searchGoods.NewBookPrice, &searchGoods.OldBookAmount, &searchGoods.OldBookPrice, &searchGoods.CreateAt, &searchGoods.UpdateAt, &searchGoods.IsSelling, &hasNewBook, &hasOldBook)
		if err != nil {
			return nil, err
		}
		if hasNewBook {
			newbookModel = &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 0, Price: searchGoods.NewBookPrice, Amount: searchGoods.NewBookAmount}
		}
		if hasOldBook {
			oldbookModel = &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 1, Price: searchGoods.OldBookPrice, Amount: searchGoods.OldBookAmount}
		}

		r = append(r, &pb.GoodsSearchResult{Book: book, GoodsId: searchGoods.GetId(), StoreId: searchGoods.StoreId, UpdateAt: searchGoods.UpdateAt, NewBook: newbookModel, OldBook: oldbookModel})

	}
	return r, nil
}

//SearchGoodsLoaction 搜索图书的货架位
func SearchGoodsLoaction(goods_id string, searchType int) (l []*pb.GoodsLocation, err error) {
	query := "select id,goods_id,type,storehouse_id,shelf_id,floor_id,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer from goods_location where 1=1"
	if searchType != -100 {
		if searchType == 0 {
			query += fmt.Sprintf(" and type=%d ", 0)
		} else {
			query += fmt.Sprintf(" and type=%d ", 1)
		}
	}
	query += fmt.Sprintf(" and goods_id='%s' order by id", goods_id)

	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return l, nil
	}
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		location := &pb.GoodsLocation{}
		l = append(l, location)

		err = rows.Scan(&location.Id, &location.GoodsId, &location.Type, &location.StorehouseId, &location.ShelfId, &location.FloorId, &location.CreateAt, &location.UpdateAt)
		if err != nil {
			return nil, err
		}
	}
	return l, err
}

//获取图书信息 精确搜索
func GetGoodsByIdOrIsbn(goods *pb.Goods) error {
	query := "select id,book_id,store_id,isbn,new_book_amount,old_book_amount,new_book_price,old_book_price,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer,is_selling,has_new_book,has_old_book from goods where 1=1"

	var args []interface{}
	var condition string

	if goods.Id != "" {
		args = append(args, goods.Id)
		condition += fmt.Sprintf(" and id=$%d", len(args))
	}

	if goods.Isbn != "" {
		args = append(args, goods.Isbn)
		condition += fmt.Sprintf(" and isbn=$%d", len(args))
	}
	args = append(args, goods.StoreId)
	condition += fmt.Sprintf(" and store_id=$%d limit 1", len(args))

	query += condition
	log.Debugf(query+"%#v", args)
	err := DB.QueryRow(query, args...).Scan(&goods.Id, &goods.BookId, &goods.StoreId, &goods.Isbn, &goods.NewBookAmount, &goods.OldBookAmount, &goods.NewBookPrice, &goods.OldBookPrice, &goods.CreateAt, &goods.UpdateAt, &goods.IsSelling, &goods.HasNewBook, &goods.HasOldBook)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//DelGoods 删除商品
func DelGoods(delModel *pb.DelGoodsModel) error {
	//首先删除货架位置
	query := "delete from goods_location where goods_id=$1 and type=$2"
	log.Debugf("delete goods_location where goods_id=%s and type=%d", delModel.Id, delModel.NewOrOld)
	_, err := DB.Exec(query, delModel.Id, delModel.NewOrOld)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	if delModel.NewOrOld == 0 {
		query = "update goods set has_new_book=false,update_at=now(),new_book_amount=0 where id=$1"
		log.Debugf("update goods set has_new_book=false,update_at=now(),new_book_amount=0 where id=%s", delModel.Id)
		_, err = DB.Exec(query, delModel.Id)
		if err != nil {
			misc.LogErr(err)
			return err
		}
	} else {
		query = "update goods set has_old_book=false,update_at=now(),old_book_amount=0 where id=$1"
		log.Debugf("update goods set has_old_book=false,update_at=now(),old_book_amount=0 where id=%s", delModel.Id)
		_, err = DB.Exec(query, delModel.Id)
		if err != nil {
			misc.LogErr(err)
			return err
		}
	}
	return nil
}

//DelGoods 下架商品
func RemoveGoods(delModel *pb.DelGoodsModel) error {
	if delModel.NewOrOld == 0 {
		query := "update goods set update_at=now(),new_book_amount=0 where id=$1"
		log.Debugf("update goods set update_at=now(),new_book_amount=0 where id=%s", delModel.Id)
		_, err := DB.Exec(query, delModel.Id)
		if err != nil {
			misc.LogErr(err)
			return err
		}
	} else {
		query := "update goods set update_at=now(),old_book_amount=0 where id=$1"
		log.Debugf("update goods set update_at=now(),old_book_amount=0 where id=%s", delModel.Id)
		_, err := DB.Exec(query, delModel.Id)
		if err != nil {
			misc.LogErr(err)
			return err
		}
	}
	return nil
}

//更新货架位信息
func UpdateGoodsLocation(goodsLocation *pb.GoodsLocation) error {
	//更新货架位置
	query := "update goods_location set storehouse_id=$1,shelf_id=$2,floor_id=$3 where id=$4 and goods_id=$5 "
	log.Debugf("update goods_location set storehouse_id=%s,shelf_id=%s,floor_id=%s where id=%s and goods_id=%s", goodsLocation.StorehouseId, goodsLocation.ShelfId, goodsLocation.FloorId, goodsLocation.Id, goodsLocation.GoodsId)
	_, err := DB.Exec(query, goodsLocation.StorehouseId, goodsLocation.ShelfId, goodsLocation.FloorId, goodsLocation.Id, goodsLocation.GoodsId)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//删除货架位置
func DelGoodsLocation(goodsLocation *pb.GoodsLocation) error {
	//删除货架位置
	query := "delete from goods_location where goods_id=$1 and id=$2"
	log.Debugf("delete from goods_location where goods_id=%s and id=%s", goodsLocation.GoodsId, goodsLocation.Id)
	_, err := DB.Query(query, goodsLocation.GoodsId, goodsLocation.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//InsertGoodsLocation 插入货架位置
func InsertGoodsLocation(loc *pb.GoodsLocation) error {
	//首先查找货架位
	query := "select id from goods_location where goods_id=$1 and type=$2 and storehouse_id=$3 and shelf_id=$4 and floor_id=$5"
	log.Debugf("select id from goods_location where goods_id=%s and type=%d and storehouse_id=%s and shelf_id=%s and floor_id=%s", loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId)
	err := DB.QueryRow(query, loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId).Scan(&loc.Id)
	//如果检查失败
	if err == sql.ErrNoRows {
		//如果用户没有上传过货架位

		query = "insert into goods_location (goods_id,type,storehouse_id,shelf_id,floor_id) values($1,$2,$3,$4,$5) returning id"

		log.Debugf("insert into goods_location (goods_id,type,storehouse_id,shelf_id,floor_id) values(%s,%d,%s,%s,%s) returning id", loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId)
		err = DB.QueryRow(query, loc.GoodsId, loc.Type, loc.StorehouseId, loc.ShelfId, loc.FloorId).Scan(&loc.Id)
		if err != nil {
			log.Errorf("%+v", err)
			return err
		}
	} else if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	return nil
}

//recover goods amount
func RecoverGoodsAmountFromClosedOrder(tx *sql.Tx, orderitem *pb.OrderItem) error {
	query := "update goods set %s where id='%s'"
	var condition string
	//新书
	if orderitem.Type == 0 {
		condition = fmt.Sprintf("new_book_amount=new_book_amount+%d", orderitem.Amount)
	} else {
		//二手书
		condition = fmt.Sprintf("old_book_amount=old_book_amount+%d", orderitem.Amount)
	}

	query = fmt.Sprintf(query, condition, orderitem.GoodsId)
	_, err := tx.Exec(query)
	if err != nil {
		log.Error(err)
		misc.LogErr(err)
		return err
	}

	return nil
}

//获取商品详情
func GetGoodsLocationDetailByIdAndType(goodsId string, bookType int64) (locations []*pb.GoodsLocationWithName, err error) {
	query := "select goods_id,storehouse_id,shelf_id,floor_id from goods_location where 1=1"

	if bookType == 0 {
		query += fmt.Sprintf(" and type=%d ", bookType)
	} else {
		query += fmt.Sprintf(" and type=%d ", bookType)
	}

	query += fmt.Sprintf(" and goods_id='%s' order by id", goodsId)

	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return locations, nil
	}
	if err != nil {

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		//遍历location
		location := &pb.GoodsLocationWithName{}
		locations = append(locations, location)
		err = rows.Scan(&location.GoodsId, &location.StorehouseId, &location.ShelfId, &location.FloorId)
		if err != nil {
			return nil, err
		}
		var id, name string
		query = fmt.Sprintf("select id,name from location where id in ('%s','%s','%s')", location.StorehouseId, location.ShelfId, location.FloorId)
		log.Debug(query)
		rows, err = DB.Query(query)
		if err != nil && err == sql.ErrNoRows {
			return locations, nil
		} else if err != nil {
			log.Error(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Error(err)
				return
			}
			if id == location.StorehouseId {
				location.StorehouseName = name
			} else if id == location.ShelfId {
				location.ShelfName = name
			} else {
				location.FloorName = name
			}
		}

	}
	return locations, err
}
