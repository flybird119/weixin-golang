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
	query := "insert into circular (store_id,type,image) values($1,$2,$3)"
	log.Debugf("insert into circular (store_id,type,image) values(%s,%d,%s)", circular.StoreId, circular.Type, circular.Image)

	_, err := DB.Exec(query, circular.StoreId, circular.Type, circular.Image)
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
	query := "update circular set update_at=now()"

	var args []interface{}
	var condition string

	if circular.Type != 0 {
		if circular.Type == 1 {
			condition += fmt.Sprintf(",type=1,url='',source_id=''")
		} else {
			if circular.SourceId != "" {
				args = append(args, circular.Type)
				condition += fmt.Sprintf(",type=$%d,source_id=$%d", len(args), len(args)+1)
				args = append(args, circular.SourceId)
				//生成URL地址
			}
		}

	}

	if circular.Title != "" {
		args = append(args, circular.Title)
		condition += fmt.Sprintf(",title=$%d", len(args))
	}

	if circular.Profile != "" {
		args = append(args, circular.Profile)
		condition += fmt.Sprintf(",profile=$%d", len(args))
	}
	if circular.Image != "" {
		args = append(args, circular.Image)
		condition += fmt.Sprintf(",image=$%d", len(args))
	}
	args = append(args, circular.Id)
	condition += fmt.Sprintf(" where id=$%d", len(args))
	query += condition
	log.Debugf(query+" args:%#v", args)

	_, err := DB.Exec(query, args...)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//CircularList 轮播图列表
func CircularList(circular *pb.Circular) (circulars []*pb.Circular, err error) {
	query := "select id,store_id,type,title,profile,image,source_id,extract(epoch from update_at)::integer from circular where store_id=$1 order by id"
	log.Debugf("select id,store_id,type,title,profile,image,source_id,extract(epoch from update_at)::integer from circular where store_id=%s order by id", circular.StoreId)
	rows, err := DB.Query(query, circular.StoreId)

	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		next := &pb.Circular{}
		circulars = append(circulars, next)
		err = rows.Scan(&next.Id, &next.StoreId, &next.Type, &next.Title, &next.Profile, &next.Image, &next.SourceId, &next.UpdateAt)
		if err != nil {
			misc.LogErr(err)
			return nil, err
		}
	}
	return circulars, nil
}
