package db

import (
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//SaveSchool 保存学校
func SaveSchool(school *pb.School) error {
	query := "insert into school (name,tel,express_fee,store_id,lat,lng) values($1,$2,$3,$4,$5,$6) returning id,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer"
	log.Debugf("insert into school (name,tel,express_fee,store_id,lat,lng) values( %s,%s,%d,%s,%f,%f ) returning id", school.Name, school.Tel, school.ExpressFee, school.StoreId, school.Lat, school.Lng)
	err := DB.QueryRow(query, school.Name, school.Tel, school.ExpressFee, school.StoreId, school.Lat, school.Lng).Scan(&school.Id, &school.CreateAt, &school.CreateAt)
	if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	return nil
}

//UpdateSchool update
func UpdateSchool(school *pb.School) error {
	updateTime := time.Now()
	query := "update school set name=$1,tel=$2,express_fee=$3,lat=$4,lng=$5,update_at=$6 where id=$7"
	log.Debugf("update school set name=$1,tel=$2,express_fee=$3,lat=$4,lng=$5 ,update_at=$6,where id=$6", school.Name, school.Tel, school.ExpressFee, school.Lat, school.Lng, updateTime, school.Id)
	_, err := DB.Query(query, school.Name, school.Tel, school.ExpressFee, school.Lat, school.Lng, updateTime, school.Id)
	if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	return nil
}

func UpdateExpressFee(school *pb.School) error {
	query := "update school set express_fee=$1 where id=$2"
	log.Debugf("update school set express_fee=%d where id=%s", school.ExpressFee, school.Id)
	_, err := DB.Exec(query, school.ExpressFee, school.Id)
	if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	return nil
}

//GetSchoolsByStore根据店铺获取所管理的学校
func GetSchoolsByStore(storeId string) (s []*pb.School, err error) {
	query := "select id, name ,tel,express_fee,lat,lng,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer from school where store_id=$1 order by create_at"
	log.Debugf("select id, name ,tel,express_fee,lat,lng,extract(epoch from create_at)::integer,extract(epoch from update_at)::integer from school where store_id=%s order by create_at", storeId)
	rows, err := DB.Query(query, storeId)

	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var school pb.School
		s = append(s, &school)
		err = rows.Scan(&school.Id, &school.Name, &school.Tel, &school.ExpressFee, &school.Lat, &school.Lng, &school.CreateAt, &school.UpdateAt)
		if err != nil {
			log.Errorf("%+v", err)
		}
	}
	if err = rows.Err(); err != nil {
		log.Debug("scan rows err last error: %s", err)
		return nil, err
	}
	return s, nil
}
