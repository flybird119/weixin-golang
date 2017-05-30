package db

import (
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
)

func DeleteLocation(loc *pb.Location) error {
	err := GetChildLocations(loc)
	if err != nil {
		log.Error(err)
		return err
	}

	query := "delete from location where id = $1"
	_, err = DB.Exec(query, loc.Id)
	log.Debugf("delete from location where id = '%s'", loc.Id)

	if err != nil {
		log.Error(err)
		return err
	}

	// delete it's children
	if len(loc.Children) > 0 {
		for _, tmp := range loc.Children {
			err = DeleteLocation(tmp)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}

// 获取子孙节点位置
func GetDescLocation(loc *pb.Location, genaration int64) error {
	err := GetChildLocations(loc)
	if err != nil {
		log.Error(err)
		return err
	}
	genaration--
	log.Debugf("The genaration is %d", genaration)

	// 默认只取出一层
	if genaration > 0 {
		// 对每个子位置取元素
		for _, loc := range loc.Children {
			err := GetDescLocation(loc, genaration)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}

	return nil
}

// 获取子节点位置
func GetChildLocations(loc *pb.Location) error {
	query := "select id, level, pid, store_id, name, extract(epoch from create_at)::integer create_at, extract(epoch from update_at)::integer update_at from location where pid = $1 and store_id = $2 order by create_at ASC"

	rows, err := DB.Query(query, loc.Id, loc.StoreId)
	log.Debugf("select id, level, pid, store_id, name, extract(epoch from create_at)::integer create_at, extract(epoch from update_at)::integer update_at from location where pid = '%s' and store_id = '%s' order by create_at ASC", loc.Id, loc.StoreId)

	if err != nil {
		log.Error(err)
		return err
	}

	for rows.Next() {
		tmp := &pb.Location{}
		err := rows.Scan(&tmp.Id, &tmp.Level, &tmp.Pid, &tmp.StoreId, &tmp.Name, &tmp.CreateAt, &tmp.UpdateAt)
		if err != nil {
			log.Error(err)
			return err
		}
		loc.Children = append(loc.Children, tmp)
		log.JSON("%+v", tmp)
	}

	return nil
}

// 单纯罗列位置
func ListLocation(loc *pb.Location) ([]*pb.Location, error) {
	query := "select id, level, pid, store_id, name, extract(epoch from create_at)::integer create_at, extract(epoch from update_at)::integer update_at from location where store_id = $1 %s order by create_at ASC"

	conditions := ""
	conditions += fmt.Sprintf("and level = %d", loc.Level)

	if loc.Pid != "" {
		// constraint pid
		conditions += fmt.Sprintf("and pid = '%s'", loc.Pid)
	}

	rows, err := DB.Query(fmt.Sprintf(query, conditions), loc.StoreId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	locations := []*pb.Location{}
	for rows.Next() {
		tempLoc := &pb.Location{}
		err = rows.Scan(&tempLoc.Id, &tempLoc.Level, &tempLoc.Pid, &tempLoc.StoreId, &tempLoc.Name, &tempLoc.CreateAt, &tempLoc.UpdateAt)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		// 获取子节点数量
		childrenAmount, err := getChildrenAmount(tempLoc.Id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		tempLoc.ChildrenAmount = childrenAmount

		locations = append(locations, tempLoc)
	}

	return locations, nil
}

func getChildrenAmount(pid string) (int64, error) {
	var total int64
	query := "select count(*) from location where pid = $1"
	log.Debugf("select count(*) from location where pid = '%s'", pid)

	err := DB.QueryRow(query, pid).Scan(&total)
	if err != nil {
		log.Error(err)
		return total, err
	}

	return total, nil
}

func UpdateLocation(loc *pb.Location) error {
	query := "update location set name = $1 where id = $2"

	_, err := DB.Exec(query, loc.Name, loc.Id)
	log.Debugf("update location set name = '%s' where id = '%s'", loc.Name, loc.Id)

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func AddLocation(location *pb.Location) error {
	query := "insert into location(level, pid, store_id, name) values($1, $2, $3, $4) returning id"

	log.Debugf("insert into location(level, pid, store_id, name) values(%d, '%s', '%s', '%s') returning id", location.Level, location.Pid, location.StoreId, location.Name)

	err := DB.QueryRow(query, location.Level, location.Pid, location.StoreId, location.Name).Scan(&location.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
