package db

import (
	"database/sql"
	"errors"
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"

	"github.com/wothing/log"
)

//通用专业批量增加
func SaveMarjor(tx *sql.Tx, major *pb.SharedMajor) error {
	query := "insert into shared_major (no,name) select '%s','%s' where not exists (select * from shared_major where name='%s') "
	query = fmt.Sprintf(query, major.No, major.Name, major.Name)
	log.Debug(query)
	_, err := tx.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//获取专业列表（筛选获取）
func FindMajorList(major *pb.SharedMajor) (models []*pb.SharedMajor, err error, totalCount int64) {
	query := "select id,no,name,extract(epoch from create_at)::bigint from shared_major where 1=1"
	queryCount := "select count(*) from shared_major where 1=1"
	var condition string
	if major.Name != "" {
		condition += fmt.Sprintf(" and name like '%s'", misc.FazzyQuery(major.Name))
	}
	//查询数量
	queryCount += condition
	log.Debug(queryCount)
	err = DB.QueryRow(queryCount).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}
	if totalCount <= 0 {
		return
	}
	if major.Page <= 0 {
		major.Page = 1
	}
	if major.Size <= 0 {
		major.Size = 15
	}
	condition += fmt.Sprintf(" order by id desc limit %d offset %d", major.Size, major.Size*(major.Page-1))
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		model := &pb.SharedMajor{}
		models = append(models, model)
		err = rows.Scan(&model.Id, &model.No, &model.Name, &model.CreateAt)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

//创建学校的学院
func SaveSchoolInstitute(model *pb.SchoolInstitute) error {
	query := "insert into map_school_institute(school_id,name) select '%s','%s' where not exists (select * from map_school_institute where school_id='%s' and name='%s') returning id"
	query = fmt.Sprintf(query, model.SchoolId, model.Name, model.SchoolId, model.Name)
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&model.Id)
	if err == sql.ErrNoRows {
		return errors.New("你已添加过此学院，请勿重复添加")
	} else if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//创建学院专业
func SaveInstituteMajor(model *pb.InstituteMajor) error {
	query := "insert into map_institute_major(institute_id,name) select '%s','%s' where not exists (select * from map_institute_major where institute_id='%s' and name='%s') returning id"
	query = fmt.Sprintf(query, model.InstituteId, model.Name, model.InstituteId, model.Name)
	log.Debug(query)

	err := DB.QueryRow(query).Scan(&model.Id)
	if err == sql.ErrNoRows {
		return errors.New("你已添加过此专业，请勿重复添加")
	} else if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//获取学校学院专业列表
func GetSchoolMajorInfo(model *pb.SchoolMajorInfoReq) (schools []*pb.GrouponSchool, err error) {
	query := "select s.id,s.name,si.id,si.name,extract(epoch from si.create_at)::bigint,im.id,im.name,extract(epoch from im.create_at)::bigint from school s left join map_school_institute si on s.id=si.school_id::uuid left join map_institute_major im on si.id=im.institute_id where 1=1"
	condition := fmt.Sprintf(" and s.status=0 and s.del_at is null and store_id='%s'", model.StoreId)

	if model.SchoolId != "" {
		condition += fmt.Sprintf(" and s.id='%s'", model.StoreId)
	}
	if model.InstituteId != "" {
		condition += fmt.Sprintf(" and si.id='%s'", model.InstituteId)
	}
	condition += " order by im.id desc"
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	//数据处理
	for rows.Next() {
		school := &pb.GrouponSchool{}
		institute := &pb.SchoolInstitute{}
		major := &pb.InstituteMajor{}

		err = rows.Scan(&school.Id, &school.Name, &institute.Id, &institute.Name, &institute.CreateAt, &major.Id, &major.Name, &major.CreateAt)
		if err != nil {
			return
		}
		var find bool
		//构建结构
		for i := 0; i < len(schools); i++ {
			if schools[i].Id == school.Id {
				institutes := schools[i].Institutes
				for j := 0; j < len(institutes); j++ {
					if institutes[j].Id == institute.Id {
						institutes[j].Majors = append(institutes[j].Majors, major)
						find = true
						break
					}
				}
				if !find {
					find = true
					institute.Majors = append(institute.Majors, major)
					schools[i].Institutes = append(schools[i].Institutes, institute)
					break
				}

			}
		}
		if !find {
			institute.Majors = append(institute.Majors, major)
			school.Institutes = append(school.Institutes, institute)
			schools = append(schools, school)
		}

	}
	return
}
