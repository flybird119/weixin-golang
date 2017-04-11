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
	//首先根据isbn获取当前用户有没有保存goods
	query := "select id from goods where isbn=$1"
	log.Debugf("select id from goods where isbn=%s", goods.Isbn)
	err := DB.QueryRow(query, goods.Isbn).Scan(&goods.Id)
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
	}
	//遍历location
	for i := 0; i < len(goods.Location); i++ {
		goods.Location[i].GoodsId = goods.Id
		err = AddGoodsLoaction(goods.Location[i])
		if err != nil {
			log.Errorf("%+v", err)
			return err
		}
	}
	return nil

}

//AddGoodsLoaction 增加货架位 goods_id  type storehouse_id shelf_id floor_id amount
func AddGoodsLoaction(loc *pb.GoodsLocation) error {
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
		query = query + " set new_book_amount=new_book_amount+$1,new_book_price=$2"
		debugQuery = debugQuery + " set new_book_amount=new_book_amount+%d,new_book_price=%d"
	} else if loc.Type == 1 {
		query = query + " set old_book_amount=old_book_amount+$1,old_book_price=$2"
		debugQuery = debugQuery + " set old_book_amount=old_book_amount+%d,old_book_price=%d"
	}
	updateTime := time.Now()
	//打开销售状态
	query = query + ",update_at=$3,is_selling=true"
	debugQuery = debugQuery + ",update_at=%f,is_selling=true"
	//修改时间
	log.Debugf(debugQuery, loc.Amount, loc.Price, updateTime)
	_, err = DB.Exec(query, loc.Amount, loc.Price, updateTime)
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
func SearchGoods(goods *pb.Goods) (r []*pb.GoodsSearchResult, err error) {
	query := "select %s from books b join goods g on b.id = g.book_id where 1=1 and is_selling=true "
	param := "b.id,b.store_id,b.title,b.isbn,b.price,b.author,b.publisher,b.pubdate,b.subtitle,b.image,b.summary,g.id, g.store_id,g.new_book_amount,g.new_book_price,g.old_book_amount,g.old_book_price,extract(epoch from g.create_at)::integer,extract(epoch from g.update_at)::integer,g.is_selling"
	query = fmt.Sprintf(query, param)
	//动态拼接参数
	var args []interface{}
	var condition string

	if goods.Isbn != "" {
		args = append(args, goods.Isbn)
		condition += fmt.Sprintf(" and b.isbn=$%d", len(args))
	}
	if goods.SearchAmount != -100 {
		args = append(args, goods.SearchAmount)
		condition += fmt.Sprintf(" and (g.new_book_amount=$%d or g.old_book_amount=$%d)", len(args), len(args)+1)
		args = append(args, goods.SearchAmount)
	}

	if goods.SearchType != -100 {
		if goods.SearchType == 0 {
			condition += " and exists (select * from goods_location gl where gl.goods_id=g.id and type =0)"
		} else {
			condition += " and exists (select * from goods_location gl where gl.goods_id=g.id and type =1)"
		}
	} else {
		condition += " and exists (select * from goods_location gl where gl.goods_id=g.id)"
	}
	if goods.SearchPicture != -100 {
		if goods.SearchPicture == 0 {
			condition += " and b.image<>''"
		} else {
			condition += " and b.image=''"
		}
	}

	if goods.Page <= 0 {
		goods.Page = 1
	}
	if goods.Size <= 0 {
		goods.Size = 20
	}
	if goods.Title != "" {
		args = append(args, misc.FazzyQuery(goods.Title))
		condition += fmt.Sprintf(" and b.title like $%d", len(args))
		args = append(args, goods.Title)
		condition += fmt.Sprintf(" order by  title <-> '$%d' ,g.id", len(args))
	} else {
		condition += " order by g.id"
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

		var salesModels []*pb.GoodsSalesModel

		/**	param := "b.id,b.store_id,b.title,b.isbn,b.price,b.author,b.publisher,b.pubdate,b.subtitle,b.image,b.summary,g.id, g.store_id,g.new_book_amount,g.new_book_price,g.old_book_amount,g.old_book_price,extract(epoch from g.create_at)::integer,extract(epoch from g.update_at)::integer,g.is_selling"
		 */
		//遍历数据
		err = rows.Scan(&book.Id, &book.StoreId, &book.Title, &book.Isbn, &book.Price, &book.Author, &book.Publisher, &book.Pubdate, &book.Subtitle, &book.Image, &book.Summary, &searchGoods.Id, &searchGoods.StoreId, &searchGoods.NewBookAmount, &searchGoods.NewBookPrice, &searchGoods.OldBookAmount, &searchGoods.OldBookPrice, &searchGoods.CreateAt, &searchGoods.UpdateAt, &searchGoods.IsSelling)
		if err != nil {
			return nil, err
		}

		if goods.SearchType == -100 {
			newLocations, _ := SearchGoodsLoaction(searchGoods.Id, 0)
			oldLocations, _ := SearchGoodsLoaction(searchGoods.Id, 1)
			log.Debugf("<<<<<<<<<<<<<<<<<<%d", goods.SearchType)
			log.Debug(newLocations == nil)
			log.Debug(oldLocations == nil)

			if newLocations != nil {
				newSales := &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 0, Price: searchGoods.NewBookPrice, Amount: searchGoods.NewBookAmount, Location: newLocations}
				salesModels = append(salesModels, newSales)
			}
			if oldLocations != nil {
				oldSales := &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 1, Price: searchGoods.OldBookPrice, Amount: searchGoods.OldBookAmount, Location: oldLocations}
				salesModels = append(salesModels, oldSales)
			}

		} else {
			if goods.SearchType == 0 {
				newLocations, _ := SearchGoodsLoaction(searchGoods.Id, 0)
				if newLocations != nil {
					newSales := &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 0, Price: searchGoods.NewBookPrice, Amount: searchGoods.NewBookAmount, Location: newLocations}
					salesModels = append(salesModels, newSales)
				}
			} else {
				oldLocations, _ := SearchGoodsLoaction(searchGoods.Id, 1)
				if oldLocations != nil {
					oldSales := &pb.GoodsSalesModel{GoodsId: searchGoods.GetId(), Type: 1, Price: searchGoods.OldBookPrice, Amount: searchGoods.OldBookAmount, Location: oldLocations}
					salesModels = append(salesModels, oldSales)
				}
			}
		}

		r = append(r, &pb.GoodsSearchResult{Book: book, Data: salesModels})
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
	if err != nil {
		log.Debugf("==========《《《==========>%+v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		location := &pb.GoodsLocation{}
		l = append(l, location)
		log.Debugf("====================>%+v", l)
		err = rows.Scan(&location.Id, &location.GoodsId, &location.Type, &location.StorehouseId, &location.ShelfId, &location.FloorId, &location.CreateAt, &location.UpdateAt)
		if err != nil {
			return nil, err
		}
	}
	return l, err
}
