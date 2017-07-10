package db

import (
	"database/sql"
	"fmt"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/pb"
	"github.com/wothing/log"
)

//保存学生学籍信息
func SaveUserSchoolStatus(model *pb.UserSchoolStatus) error {
	searchModel := &pb.UserSchoolStatus{UserId: model.UserId}
	err := GetUserSchoolStatus(searchModel)
	if err != nil {
		log.Error(err)
		return err
	}
	if searchModel.Id != "" {
		model.Id = ""
		return nil
	}
	query := "insert into user_school_status(school_id,user_id,institute_id,institute_major_id) select '%s','%s','%s','%s'  returning id"
	query = fmt.Sprintf(query, model.SchoolId, model.UserId, model.InstituteId, model.InstituteMajorId)
	log.Debug(query)
	err = DB.QueryRow(query).Scan(&model.Id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//更新学生学籍信息
func UpdateUserSchoolStatus(model *pb.UserSchoolStatus) error {
	query := "update user_school_status set update_at= now()"
	var condition string
	if model.SchoolId != "" {
		condition += fmt.Sprintf(",school_id='%s'", model.SchoolId)
	}
	if model.InstituteId != "" {
		condition += fmt.Sprintf(",institute_id='%s'", model.InstituteId)

	}
	if model.InstituteMajorId != "" {
		condition += fmt.Sprintf(",institute_major_id='%s'", model.InstituteMajorId)

	}
	condition += fmt.Sprintf(" where id='%s'", model.Id)
	query += condition
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//获取学生学籍
func GetUserSchoolStatus(model *pb.UserSchoolStatus) error {
	query := "select us.id,us.school_id,us.user_id,us.institute_id,us.institute_major_id extract(epoch from us.create_at)::bigint,s.name,si.name,im.name from user_school_status us join school s on us.school_id=s.id join map_school_institute si on s.id=si.school_id join map_institute_major im on si.id=im.institute_id"
	var condition string
	condition += " where s.status=0 and s.del_at is null and si.status=1 and im.status=1 and user_id='%s' limit 1 order by us.create_at desc"
	condition = fmt.Sprintf(condition, model.UserId)
	query += condition
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&model.Id, &model.SchoolId, &model.UserId, &model.InstituteId, &model.InstituteMajorId, &model.CreateAt, &model.SchoolName, &model.InstituteName, &model.MajorName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			log.Error(err)
			return err
		}
	}
	return nil
}
