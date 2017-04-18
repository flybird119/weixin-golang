package db

import (
	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//CartAdd 增加购物车
func CartAdd(cart *pb.Cart) error {
	//首先查找当前用户购物车又没有这个商品
	query := "select id from cart where user_id=$1 and store_id=$2 and goods_id=$3 and type=$4"
	log.Debugf("select id from cart where user_id=%s and store_id=%s and goods_id=%s and type=%d", cart.UserId, cart.StoreId, cart.GoodsId, cart.Type)
	err := DB.QueryRow(query, cart.UserId, cart.StoreId, cart.GoodsId, cart.Type).Scan(&cart.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	//判断当前用户购物车又没有这种类型的商品
	if cart.Id == 0 {
		//购物车没有改类型商品，添加
		query = "insert into cart (user_id,store_id,goods_id,type,amount) values($1,$2,$3,$4,$5)"
		log.Debugf("insert into cart (user_id,store_id,goods_id,type,amount) values(%s,%s,%s,%d,%d)", cart.UserId, cart.StoreId, cart.GoodsId, cart.Type, cart.Amount)
		_, err = DB.Exec(query, cart.UserId, cart.StoreId, cart.GoodsId, cart.Type, cart.Amount)
		if err != nil {
			misc.LogErr(err)
			return err
		}
		return nil
	} else {
		//购物车有此商品
		query = "update cart set amount=amount+$1 where id=$2"
		log.Debugf("update cart set amount=amount+%d where id=%s", cart.Amount, cart.Id)
		if err != nil {
			misc.LogErr(err)
			return err
		}
		return nil
	}
}

//CartList 购物车列表
func CartList(cart *pb.Cart) (carts []*pb.Cart, err error) {
	query := "select id,user_id,store_id,goods_id,type,amount from cart where user_id=$1 and store_id=$2"
	log.Debugf("select id,user_id,store_id,goods_id,type,amount from cart where user_id=%s and store_id=%s", cart.UserId, cart.StoreId)
	rows, err := DB.Query(query, cart.UserId, cart.StoreId)
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		nextCart := &pb.Cart{}
		carts = append(carts, nextCart)
		err = rows.Scan(&nextCart.Id, &nextCart.UserId, &nextCart.StoreId, &nextCart.GoodsId, &nextCart.Type, &nextCart.Amount)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
	}
	return carts, err
}

//CartUpdate 更改购物车
func CartUpdate(cart *pb.Cart) error {
	query := "update cart set amount=$1 where id=$2"
	log.Debugf("update cart set amount=%d where id=%s", cart.Amount, cart.Id)
	_, err := DB.Exec(query, cart.Amount, cart.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//CartDel 删除购物车
func CartDel(cart *pb.Cart) error {
	query := "delete from cart where id=$1 and user_id=$2"
	log.Debugf("delete from cart where id=%s and user_id=%s", cart.Id, cart.UserId)
	_, err := DB.Exec(query, cart.Id, cart.UserId)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}