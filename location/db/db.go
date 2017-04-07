package db

import (
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
)

func ListLocation(loc *pb.Location) ([]*pb.Location, error) {
	query := "select id, level, pid, store_id, name, extract(epoch from create_at)::integer create_at, extract(epoch from update_at)::integer update_at from location where store_id = $1 %s order by create_at DESC"

	conditions := ""
	if loc.Pid != "" {
		// constraint pid
		conditions += fmt.Sprintf("and pid = '%s'", loc.Pid)
	}

	if loc.Level != 0 {
		// constraint level
		conditions += fmt.Sprintf("and level = %d", loc.Level)
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
		locations = append(locations, tempLoc)
	}

	return locations, nil
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
	query := "insert into location(level, pid, store_id, name) values($1, $2, $3, $4)"

	log.Debugf("insert into location(level, pid, store_id, name) values(%d, '%s', '%s', '%s')", location.Level, location.Pid, location.StoreId, location.Name)

	_, err := DB.Exec(query, location.Level, location.Pid, location.StoreId, location.Name)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
