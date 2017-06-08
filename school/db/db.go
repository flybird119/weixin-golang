package db

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//SaveSchool 保存学校
func SaveSchool(school *pb.School) error {
	query := "insert into school (name,tel,express_fee,store_id,lat,lng) values($1,$2,$3,$4,$5,$6) returning id,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint"
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

//更改运费
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

//删除学校
func DelSchool(school *pb.School) error {

	var query string
	//查找这个学校线上订单销售
	orderCount, err := getOrderCountBySchool(school.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	//查找这个学校线上订单销售
	retailCount, err := getRetailCountBySchool(school.Id)

	if err != nil {
		log.Error(err)
		return err
	}

	if orderCount == 0 && retailCount == 0 {
		//物理删除学校
		query = fmt.Sprintf("delete from school where id='%s'", school.Id)
	} else {
		//逻辑 删除
		query = fmt.Sprintf("update school set status=1,del_at=now(),del_staff_id='%s',update_at=now() where id='%s'", school.DelStaffId, school.Id)
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	tx.Commit()
	return nil
}

//GetSchoolsByStore根据店铺获取所管理的学校
func GetSchoolsByStore(storeId string, status int64) (s []*pb.School, err error) {
	query := "select id, name ,tel,express_fee,lat,lng,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint,status,extract(epoch from del_at)::bigint,del_staff_id from school where 1=1"

	var condition string
	if status != 3 {
		condition += fmt.Sprintf(" and status=%d", status)
	}
	condition += fmt.Sprintf(" and store_id='%s' order by create_at", storeId)
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Errorf("%+v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var school pb.School
		s = append(s, &school)
		school.StoreId = storeId
		var delAt sql.NullInt64
		err = rows.Scan(&school.Id, &school.Name, &school.Tel, &school.ExpressFee, &school.Lat, &school.Lng, &school.CreateAt, &school.UpdateAt, &school.Status, &delAt, &school.DelStaffId)
		if err != nil {
			log.Errorf("%+v", err)
		}
		if delAt.Valid {
			school.DelAt = delAt.Int64
		}
	}
	if err = rows.Err(); err != nil {
		log.Debug("scan rows err last error: %s", err)
		return nil, err
	}
	return s, nil
}

//获取学校信息
func GetSchoolById(schoolId string) (school *pb.School, err error) {
	query := "select id, name ,tel,express_fee,lat,lng,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint,status ,extract(epoch from del_at)::bigint,del_staff_id from school where id=$1"
	log.Debugf("select id, name ,tel,express_fee,lat,lng,extract(epoch from create_at)::bigint,extract(epoch from update_at)::bigint from,status,extract(epoch from del_at)::bigint,del_staff_id school where id='%s'", schoolId)
	school = &pb.School{}
	var delAt sql.NullInt64
	err = DB.QueryRow(query, schoolId).Scan(&school.Id, &school.Name, &school.Tel, &school.ExpressFee, &school.Lat, &school.Lng, &school.CreateAt, &school.UpdateAt, &school.Status, &delAt, &school.DelStaffId)
	if err != nil {
		log.Debugf("%#v", err)
		return nil, err
	}
	if delAt.Valid {
		school.DelAt = delAt.Int64
	}

	return school, nil
}

//查看一个学校有没有线上订单
func getOrderCountBySchool(schoolId string) (int64, error) {
	var totalCount int64
	query := "select count(*) from orders where school_id='%s'"
	query = fmt.Sprintf(query, schoolId)
	err := DB.QueryRow(query).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return totalCount, nil
}

//根据学校获取线下零售条数
func getRetailCountBySchool(schoolId string) (int64, error) {
	var totalCount int64
	query := "select  count(*) from retail where school_id='%s'"
	query = fmt.Sprintf(query, schoolId)
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}
