package db

import (
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//AddCircular 增加轮播图
func AddCircular(circular *pb.Circular) error {
	query := "insert into circular (store_id,type,source_title,source_id,profile,image) values($1,$2,$3,$4,%5,$6)"
	log.Debugf("insert into circular (store_id,type,source_title,source_id,profile,image) values(%s,%d,%s,%s,%s,%s)", circular.StoreId, circular.Type, circular.SourceTitle, circular.Id, circular.Profile, circular.Image)
	_, err := DB.Exec(query, circular.StoreId, circular.Type, circular.SourceTitle, circular.Id, circular.Profile, circular.Image)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//DelCircular 删除轮播图
func DelCircular(circular *pb.Circular) error {
	query := "delete from circular where id=$1"
	log.Debugf("delete from circular where id=%s", circular.Id)
	_, err := DB.Exec(query, circular)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//UpdateCircular 更新轮播图
func UpdateCircular(circular *pb.Circular) error {
	query := "update cicular set update_at=now()"

	var args []interface{}
	var condition string

	if circular.Type != 0 {
		args = append(args, circular.Type)
		condition += fmt.Sprintf(",type=$%d", len(args))
	}

	if circular.SourceTitle != "" {
		args = append(args, circular.SourceTitle)
		condition += fmt.Sprintf(",source_title=$%d", len(args))
	}

	if circular.SourceId != "" {

	}
	query += condition
	log.Debugf(query+" args:%#v", args)

	return nil
}

//CircularList 轮播图列表
func CircularList(circular *pb.Circular) (circulars []*pb.Circular, err error) {

	return nil, nil
}
