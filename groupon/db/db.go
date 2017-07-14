package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	. "github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/misc"
	"github.com/goushuyun/weixin-golang/pb"
	schoolDB "github.com/goushuyun/weixin-golang/school/db"

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
	query := "insert into map_school_institute(school_id,name) select '%s','%s' where not exists (select * from map_school_institute where school_id='%s' and name='%s' and status=1) returning id"
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
	query := "insert into map_institute_major(institute_id,name) select '%s','%s' where not exists (select * from map_institute_major where institute_id='%s' and name='%s' and status=1) returning id"
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
	query := "select s.id,s.name,si.id,si.name,extract(epoch from si.create_at)::bigint,si.status,im.id,im.name,extract(epoch from im.create_at)::bigint,im.status from school s left join map_school_institute si on s.id=si.school_id::uuid left join map_institute_major im on si.id=im.institute_id where 1=1"
	condition := fmt.Sprintf(" and s.status=0 and s.del_at is null and s.store_id='%s'", model.StoreId)

	if model.SchoolId != "" {
		condition += fmt.Sprintf(" and s.id='%s'::uuid", model.SchoolId)
	}
	if model.InstituteId != "" {
		condition += fmt.Sprintf(" and si.id='%s'", model.InstituteId)
	}
	if model.UserType != 1 {
		condition += fmt.Sprintf(" and (si.status=1) and (im.status=1)")
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
		var iId, iName, mId, mName sql.NullString
		var iCreateAt, iStatus, mCreateAt, mStatus sql.NullInt64
		err = rows.Scan(&school.Id, &school.Name, &iId, &iName, &iCreateAt, &iStatus, &mId, &mName, &mCreateAt, &mStatus)
		if err != nil {
			return
		}
		if iId.Valid {
			institute.Id = iId.String
		}
		if iName.Valid {
			institute.Name = iName.String
		}
		if mId.Valid {
			major.Id = mId.String
		}
		if mName.Valid {
			major.Name = mName.String
		}
		if iCreateAt.Valid {
			institute.CreateAt = iCreateAt.Int64
		}
		if iStatus.Valid {
			institute.Status = iStatus.Int64
		}
		if mCreateAt.Valid {
			major.CreateAt = mCreateAt.Int64
		}
		if mStatus.Valid {
			major.Status = mStatus.Int64
		}
		var find bool

		//构建结构
		for i := 0; i < len(schools); i++ {

			if schools[i].Id == school.Id {
				institutes := schools[i].Institutes
				for j := 0; j < len(institutes); j++ {
					if institutes[j].Id == institute.Id {
						if major.Id != "" && major.Status != 2 {
							institutes[j].Majors = append(institutes[j].Majors, major)
						}
						find = true
						break

					}
				}
				if !find {
					if major.Id != "" && major.Status != 2 {
						institute.Majors = append(institute.Majors, major)
					}
					if institute.Id != "" && institute.Status != 2 {
						schools[i].Institutes = append(schools[i].Institutes, institute)
					}
					find = true
					break
				}
			}
		}
		if !find {
			if major.Id != "" && major.Status != 2 {
				institute.Majors = append(institute.Majors, major)
			}
			if institute.Id != "" && institute.Status != 2 {
				school.Institutes = append(school.Institutes, institute)
			}
			schools = append(schools, school)
		}

	}
	return
}

//创建班级购
func SaveGroupon(model *pb.Groupon) error {
	var avatar sql.NullString
	var avatarString string
	var query string
	if model.FounderType == 1 {
		query = fmt.Sprintf("select avatar from users where id='%s'", model.FounderId)
	} else {
		query = fmt.Sprintf("select logo from store where id='%s'", model.StoreId)
	}
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&avatar)
	if err != nil {
		log.Error(err)
		return err
	}
	if avatar.Valid {
		avatarString = avatar.String
	}
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer tx.Rollback()

	query = "insert into groupon(store_id,school_id,institute_id,institute_major_id,founder_id,term,class,founder_type,founder_name,founder_mobile,profile,expire_at,founder_avatar,participate_num) values(%s) returning id"
	param := fmt.Sprintf("'%s','%s','%s','%s','%s','%s','%s',%d,'%s','%s','%s',to_timestamp(%d),'%s',1", model.StoreId, model.SchoolId, model.InstituteId, model.InstituteMajorId, model.FounderId, model.Term, model.Class, model.FounderType, model.FounderName, model.FounderMobile, model.Profile, model.ExpireAt, avatarString)
	query = fmt.Sprintf(query, param)
	log.Debug(query)
	err = tx.QueryRow(query).Scan(&model.Id)
	if err != nil {
		log.Error(err)
		return err
	}
	for i := 0; i < len(model.Items); i++ {
		item := model.Items[i]
		query = fmt.Sprintf("insert into groupon_item (groupon_id,goods_id) select '%s','%s' where not exists (select * from groupon_item where groupon_id='%s' and goods_id='%s')", model.Id, item.GoodsId, model.Id, item.GoodsId)
		_, err = tx.Exec(query)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	tx.Commit()
	oplog := &pb.GrouponOperateLog{GrouponId: model.Id, FounderId: model.FounderId, FounderName: model.FounderName, FounderType: model.FounderType, OperateType: "create", OperateDetail: " "}
	err = SaveGrouponOperateLog(oplog)
	if err != nil {
		log.Error(err)
	}
	return nil
}

//班级购列表
func FindGroupon(model *pb.Groupon) (models []*pb.Groupon, err error, totalCount int64) {
	query := "select %s from groupon g join school s on s.id=g.school_id::uuid left join map_school_institute si on s.id=si.school_id::uuid left join map_institute_major im on si.id=im.institute_id where 1=1"
	param := " distinct g.id,g.status,g.store_id,g.school_id,g.institute_id,g.institute_major_id,g.founder_id,g.term,g.class,g.founder_type,g.founder_name,g.founder_mobile,g.profile,g.participate_num,g.star_num,g.total_sales,g.order_num,extract(epoch from g.create_at)::bigint,extract(epoch from g.expire_at)::bigint,s.name,s.status,si.name,si.status,im.name,im.status,g.founder_avatar"
	queryCount := fmt.Sprintf(query, "count(*)")
	query = fmt.Sprintf(query, param)
	var condition string
	condition += " and g.institute_id=si.id and g.institute_major_id=im.id"
	// 根据编号
	if model.Id != "" {
		condition += fmt.Sprintf(" and g.id='%s'", model.Id)
	}
	// 根据学校
	if model.SchoolId != "" {
		condition += fmt.Sprintf(" and g.school_id='%s'", model.SchoolId)
	}
	// 根据创建类型
	if model.FounderType != 0 {
		condition += fmt.Sprintf(" and g.founder_type=%d", model.FounderType)
	}
	// 根据学院
	if model.InstituteId != "" {
		condition += fmt.Sprintf(" and g.institute_id='%s'", model.InstituteId)
	}
	// 根据专业
	if model.InstituteMajorId != "" {
		condition += fmt.Sprintf(" and g.institute_major_id='%s'", model.InstituteMajorId)
	}
	// 根据班级
	if model.Class != "" {
		condition += fmt.Sprintf(" and g.class like '%s'", misc.FazzyQuery(model.Class))
	}
	// 根据学期
	if model.Term != "" {
		condition += fmt.Sprintf(" and g.term='%s'", model.Term)
	}
	// 根据是否可用
	if model.SearchType != 0 {
		if model.SearchType == 1 {
			//正常 使用中的
			condition += fmt.Sprintf(" and g.status=1 and g.expire_at>to_timestamp(%d)", time.Now().Unix())

		} else if model.SearchType == 2 {
			//过期
			condition += fmt.Sprintf(" and g.expire_a<to_timestamp(%d)", time.Now().Unix())

		} else if model.SearchType == 3 {
			//异常
			condition += " and g.status=2"
		}
	}
	// 根据参与者
	if model.ParticipateUser != "" {
		condition += fmt.Sprintf(" and exists (select * from groupon_operate_log where founder_id='%s' and ((operate_type='purchase') or (operate_type='share')))", model.ParticipateUser)
	}
	// 根据创建人 创建人类型
	if model.FounderId != "" {
		condition += fmt.Sprintf(" and g.founder_id='%s'", model.FounderId)
	}

	condition += fmt.Sprintf(" and g.store_id='%s'", model.StoreId)
	queryCount += condition
	log.Debug(queryCount)
	err = DB.QueryRow(queryCount).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}
	if totalCount == 0 {
		return
	}
	if model.Page <= 0 {
		model.Page = 1
	}
	if model.Size <= 0 {
		model.Size = 15
	}
	condition += fmt.Sprintf(" order by order_num desc, id desc limit %d offset %d", model.Size, model.Size*(model.Page-1))
	query += condition
	log.Debug(query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		major := &pb.InstituteMajor{}
		institute := &pb.SchoolInstitute{}
		school := &pb.GrouponSchool{}
		m := &pb.Groupon{Major: major, Institute: institute, School: school}
		models = append(models, m)
		//param := "g.id,g.status,g.store_id,g.school_id,g.institute_id,g.institute_major_id,g.founder_id,g.term,g.class,g.founder_type,g.founder_name,g.founder_mobile,g.profile,g.participate_num,g.star_num,g.total_sales,g.order_num,extract(epoch from g.create_at)::bigint,extract(epoch from g.expire_at)::bigint,s.name,s.status,si.name,si.status,im.name,im.status"
		err = rows.Scan(&m.Id, &m.Status, &m.StoreId, &m.SchoolId, &m.InstituteId, &m.InstituteMajorId, &m.FounderId, &m.Term, &m.Class, &m.FounderType, &m.FounderName, &m.FounderMobile, &m.Profile, &m.ParticipateNum, &m.StarNum, &m.TotalSales, &m.OrderNum, &m.CreateAt, &m.ExpireAt, &school.Name, &school.Status, &institute.Name, &institute.Status, &major.Name, &major.Status, &m.FounderAvatar)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

//新增班级购项
func GetGrouponItems(model *pb.Groupon) (models []*pb.GrouponItem, err error) {
	//获取图书信息
	query := "select gi.id,gi.goods_id,b.isbn,b.title,b.author,b.image,g.new_book_amount,g.old_book_amount,g.new_book_price,g.old_book_price,g.has_new_book,g.has_old_book from groupon_item gi join goods g on gi.goods_id=g.id join books b on g.book_id=b.id where gi.groupon_id='%s'"
	query = fmt.Sprintf(query, model.Id)
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &pb.GrouponItem{}
		models = append(models, m)
		err = rows.Scan(&m.Id, &m.GoodsId, &m.BookIsbn, &m.BookTitle, &m.BookAuthor, &m.BookImage, &m.NewBookAmount, &m.OldBookAmount, &m.NewBookPrice, &m.OldBookPrice, &m.HasNewBook, &m.HasOldBook)
		if err != nil {
			log.Error(err)
			return
		}

	}
	return
}

//获取班级购参与人信息
func GetGrouponPurchaseUsers(model *pb.Groupon) (models []*pb.GrouponUserInfo, err error) {
	query := "select id,nickname,avatar from users where id in (select founder_id from groupon_operate_log where groupon_id='%s' and operate_type='purchase')"
	query = fmt.Sprintf(query, model.Id)
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &pb.GrouponUserInfo{}
		models = append(models, m)
		err = rows.Scan(&m.Id, &m.Name, &m.Avatar)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

//获取班级购操作日志
func GetGrouponOperateLog(model *pb.Groupon) (models []*pb.GrouponOperateLog, err error) {
	query := "select id,founder_id,founder_type,founder_name,operate_type,operate_detail,extract(epoch from create_at)::bigint,founder_avatar from groupon_operate_log where groupon_id='%s' and operate_type <> 'star' order by id desc"
	query = fmt.Sprintf(query, model.Id)
	log.Debug(query)
	rows, err := DB.Query(query)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		log.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		m := &pb.GrouponOperateLog{}
		models = append(models, m)
		err = rows.Scan(&m.Id, &m.FounderId, &m.FounderType, &m.FounderName, &m.OperateType, &m.OperateDetail, &m.CreateAt, &m.FounderAvatar)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

//修改班级购
func UpdateGruopon(model *pb.Groupon) error {
	query := "update groupon set update_at=now()"
	var condition string
	if model.Status != 0 {
		condition += fmt.Sprintf(",status=%d", model.Status)
	}
	if model.InstituteId != "" {
		condition += fmt.Sprintf(",institute_id='%s'", model.InstituteId)
	}
	if model.InstituteMajorId != "" {
		condition += fmt.Sprintf(",institute_major_id='%s'", model.InstituteMajorId)
	}
	if model.Term != "" {
		condition += fmt.Sprintf(",term='%s'", model.Term)
	}
	if model.Class != "" {
		condition += fmt.Sprintf(",class='%s'", model.Class)
	}
	if model.FounderName != "" {
		condition += fmt.Sprintf(",founder_name='%s'", model.FounderName)
	}
	if model.FounderMobile != "" {
		condition += fmt.Sprintf(",founder_mobile='%s'", model.FounderMobile)
	}
	if model.Profile != "" {
		condition += fmt.Sprintf(",profile='%s'", model.Profile)
	}
	if model.ParticipateNum != 0 {
		condition += fmt.Sprintf(",participate_num=participate_num+%d", model.ParticipateNum)
	}
	if model.StarNum != 0 {
		condition += fmt.Sprintf(",star_num=star_num+%d", model.StarNum)
	}
	if model.TotalSales != 0 {
		condition += fmt.Sprintf(",total_sales=total_sales+%d", model.TotalSales)
	}
	if model.OrderNum != 0 {
		condition += fmt.Sprintf(",order_num=order_num+%d", model.OrderNum)
	}
	if model.ExpireAt != 0 {
		condition += fmt.Sprintf(",expire_at=to_timestamp(%d)", model.ExpireAt)
	}
	condition += fmt.Sprintf(" where id='%s' returning founder_id, founder_name,founder_type ", model.Id)
	query += condition
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&model.FounderId, &model.FounderName, &model.FounderType)
	if err != nil {
		log.Error(err)
		return err
	}
	if len(model.DelItemIds) != 0 {
		query = "delete from groupon_item where id in(%s) and groupon_id='%s'"
		stmt := strings.Repeat(",'%s'", len(model.DelItemIds))
		var ids []interface{}
		for _, value := range model.DelItemIds {
			ids = append(ids, value.Id)
		}
		condition = fmt.Sprintf(stmt, ids...)
		condition = string([]byte(condition)[1:])
		query = fmt.Sprintf(query, condition, model.Id)
		log.Debug(query)
		_, err = DB.Exec(query)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	if len(model.AddItems) != 0 {
		for i := 0; i < len(model.AddItems); i++ {
			item := model.AddItems[i]
			query = fmt.Sprintf("insert into groupon_item (groupon_id,goods_id) select '%s','%s' where not exists (select * from groupon_item where groupon_id='%s' and goods_id='%s')", model.Id, item.GoodsId, model.Id, item.GoodsId)
			_, err = DB.Exec(query)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

//保存班级购操作日志
func SaveGrouponOperateLog(model *pb.GrouponOperateLog) error {
	var avatar sql.NullString
	var avatarString string
	var query string
	if model.FounderType == 1 {
		query = fmt.Sprintf("select avatar from users where id='%s'", model.FounderId)
	} else {
		query = fmt.Sprintf("select logo from store s join groupon g on g.store_id=s.id where g.id='%s'", model.GrouponId)
	}
	log.Debug(query)
	err := DB.QueryRow(query).Scan(&avatar)
	if err != nil {
		log.Error(err)
		return err
	}
	if avatar.Valid {
		avatarString = avatar.String
	}
	query = "insert into groupon_operate_log (groupon_id,founder_id,founder_type,founder_name,operate_type,operate_detail,founder_avatar) values('%s','%s',%d,'%s','%s','%s','%s')"
	query = fmt.Sprintf(query, model.GrouponId, model.FounderId, model.FounderType, model.FounderName, model.OperateType, model.OperateDetail, avatarString)
	log.Debug(query)
	_, err = DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GrouponSubmit(orderModel *pb.GrouponSubmitModel) (order *pb.Order, noStock string, err error) {
	groupon, err := GetGrouponInfo(orderModel.GrouponId)
	if err != nil {
		log.Error(err)
		return
	}
	tx, err := DB.Begin()
	if err != nil {
		log.Error(err)
		return
	}
	defer tx.Rollback()
	//获取学校的运费
	school, err := schoolDB.GetSchoolById(groupon.SchoolId)
	if err != nil {
		misc.LogErr(err)
		return nil, "", err
	}
	nowTime := time.Now()
	order = &pb.Order{}
	//首选创建goods，然后创建订单项
	query := "insert into orders (total_fee,freight,user_id,mobile,name,address,remark,store_id,school_id,order_at,goods_fee,withdrawal_fee,groupon_id ) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,0,0,$11) returning id"
	log.Debugf(query+"args : %#v", school.ExpressFee, school.ExpressFee, orderModel.UserId, orderModel.Mobile, orderModel.Name, orderModel.Address, orderModel.Remark, orderModel.StoreId, orderModel.SchoolId, nowTime, orderModel.GrouponId)
	err = tx.QueryRow(query, school.ExpressFee, school.ExpressFee, orderModel.UserId, orderModel.Mobile, orderModel.Name, orderModel.Address, orderModel.Remark, orderModel.StoreId, orderModel.SchoolId, nowTime, orderModel.GrouponId).Scan(&order.Id)
	if err != nil {
		misc.LogErr(err)
		return nil, "", err
	}
	//遍历carts
	for i := 0; i < len(orderModel.Items); i++ {
		noStock, err = AddOrderItem(tx, order, orderModel.Items[i], nowTime)
		if err != nil {
			misc.LogErr(err)
			return nil, "", err
		}
		if noStock != "" {
			return nil, "noStock", nil
		}
	}

	query = "select order_status,total_fee,freight,goods_fee,user_id,mobile,name,address,remark,store_id,school_id,groupon_id from orders where id=$1"
	log.Debugf("select order_status,total_fee,freight,goods_fee,user_id,mobile,name,address,remark,store_id,school_id,groupon_id from orders where id='%s'", order.Id)
	err = tx.QueryRow(query, order.Id).Scan(&order.OrderStatus, &order.TotalFee, &order.Freight, &order.GoodsFee, &order.UserId, &order.Mobile, &order.Name, &order.Address, &order.Remark, &order.StoreId, &order.SchoolId, &order.GrouponId)
	if err != nil {
		misc.LogErr(err)
		return nil, "", err
	}
	tx.Commit()
	//增加团购操作日志
	oplog := &pb.GrouponOperateLog{GrouponId: orderModel.GrouponId, FounderId: orderModel.UserId, FounderName: orderModel.Name, FounderType: 1, OperateType: "purchase", OperateDetail: " "}
	err = SaveGrouponOperateLog(oplog)
	if err != nil {
		log.Error(err)
	}

	num := 1
	totalCount, _ := HasGrouponLogWithOpreation(orderModel.GrouponId, orderModel.UserId, "purchase")
	if totalCount > 1 {
		num = 0
	}
	//更新团购数据
	updateG := &pb.Groupon{Id: orderModel.GrouponId, TotalSales: order.TotalFee, ParticipateNum: int64(num), OrderNum: 1}
	err = UpdateGruopon(updateG)
	if err != nil {
		log.Error(err)
		err = nil
	}
	return
}

func AddOrderItem(tx *sql.Tx, order *pb.Order, item *pb.GrouponItem, nowTime time.Time) (noStack string, err error) {
	//减少库存量

	var (
		query      string
		price      int
		amount     int
		is_selling bool
	)
	noStock := "noStock"
	if item.Type == 0 {
		query = "update goods set new_book_amount=new_book_amount-$1,new_book_sale_amount=new_book_sale_amount+$2 where id=$3 returning new_book_amount,new_book_price,has_new_book"
		log.Debugf("update goods set new_book_amount=new_book_amount-%d,new_book_sale_amount=new_book_sale_amount+%d where id='%s'returning new_book_amount,new_book_price,has_new_book", item.Amount, item.Amount, item.GoodsId)
		err = tx.QueryRow(query, item.Amount, item.Amount, item.GoodsId).Scan(&amount, &price, &is_selling)
		if err != nil {
			misc.LogErr(err)
			return "", err
		}
		if !is_selling || amount < 0 {

			return noStock, nil
		}
	} else {
		query = "update goods set old_book_amount=old_book_amount-$1,old_book_sale_amount=old_book_sale_amount+$2 where id=$3 returning old_book_amount,old_book_price,has_old_book"
		log.Debugf("update goods set old_book_amount=old_book_amount-%d ,old_book_sale_amount=old_book_sale_amount+%s where id='%s'returning old_book_amount,old_book_price,has_old_book", item.Amount, item.GoodsId)
		err = tx.QueryRow(query, item.Amount, item.Amount, item.GoodsId).Scan(&amount, &price, &is_selling)
		if err != nil {
			misc.LogErr(err)
			return "", err
		}
		if !is_selling || amount < 0 {

			return noStock, nil
		}
	}
	//然后创建订单项
	query = "insert into orders_item (goods_id,orders_id,type,amount,price,create_at) values($1,$2,$3,$4,$5,$6)"
	log.Debugf("insert into orders_item (goods_id,orders_id,type,amount,price,create_at) values('%s','%s',%d,%d,%d,%v)", item.GoodsId, order.Id, item.Type, item.Amount, price, nowTime)
	_, err = tx.Exec(query, item.GoodsId, order.Id, item.Type, item.Amount, price, nowTime)
	if err != nil {
		misc.LogErr(err)
		return "", err
	}
	//更改订单
	totalFee := int(item.Amount) * price
	query = "update orders set total_fee=total_fee+$1,goods_fee=goods_fee+$2 where id=$3"
	log.Debugf("update orders set total_fee=total_fee+%d,goods_fee=goods_fee+%d where id='%s'", totalFee, totalFee, order.Id)
	_, err = tx.Exec(query, totalFee, totalFee, order.Id)
	if err != nil {
		misc.LogErr(err)
		return "", err
	}
	return "", nil
}

//获取团购信息
func GetGrouponInfo(id string) (groupon *pb.Groupon, err error) {
	query := "select id,status,store_id,school_id,institute_id,institute_major_id from groupon where id='%s'"
	query = fmt.Sprintf(query, id)
	log.Debug(query)
	groupon = &pb.Groupon{}
	err = DB.QueryRow(query).Scan(&groupon.Id, &groupon.Status, &groupon.StoreId, &groupon.SchoolId, &groupon.InstituteId, &groupon.InstituteMajorId)
	if err != nil {
		misc.LogErr(err)
		return
	}
	return
}

func HasGrouponLogWithOpreation(grouponId, userId, oprea string) (totalCount int64, err error) {
	query := "select count(*) from groupon_operate_log where groupon_id='%s' and founder_id='%s' and operate_type='%s'"
	query = fmt.Sprintf(query, grouponId, userId, oprea)
	log.Debug(query)

	err = DB.QueryRow(query).Scan(&totalCount)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

//删除专业
func DelInstituMajor(model *pb.InstituteMajor) error {
	query := fmt.Sprintf("update map_institute_major set status=2 where id='%s'", model.Id)
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//删除学院
func UpdateInstituteMajor(model *pb.InstituteMajor) error {
	query := fmt.Sprintf("update map_institute_major m set name='%s' where id='%s' and not exists (select * from map_institute_major m1 where m.institute_id=m1.institute_id and m1.name='%s' and m1.status=1) returning name", model.Name, model.Id, model.Name)
	log.Debug(query)
	var name sql.NullString
	err := DB.QueryRow(query).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("exists")
		}
		log.Error(err)
		return err
	}
	var vName string
	if name.Valid {
		vName = name.String
	}
	if vName == "" {
		return errors.New("exists")
	}

	return nil
}

//删除学院
func DelSchoolInstitute(model *pb.SchoolInstitute) error {
	query := fmt.Sprintf("update map_school_institute set status=2 where id='%s'", model.Id)
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//修改学校学院名称
func UpdateSchoolInstitute(model *pb.SchoolInstitute) error {
	query := fmt.Sprintf("update map_school_institute m set name='%s' where id='%s' and not exists (select * from map_school_institute m1 where m.school_id=m1.school_id and m1.name='%s' and m1.status=1)  returning name", model.Name, model.Id, model.Name)
	log.Debug(query)
	var name sql.NullString
	err := DB.QueryRow(query).Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("exists")
		}
		log.Error(err)
		return err
	}
	var vName string
	if name.Valid {
		vName = name.String
	}
	if vName == "" {
		return errors.New("exists")
	}

	return nil
}

//批量修改截止日期
func BatchUpdateGrouponExpireAt(model *pb.Groupon) error {
	query := "update groupon set expire_at=to_timestamp(%d) where id in(%s) and store_id='%s'"
	var condition string
	condition = strings.Repeat(",'%s'", len(model.UpdateIds))
	condition = condition[1:]
	var param []interface{}
	for i := 0; i < len(model.UpdateIds); i++ {
		param = append(param, model.UpdateIds[i].Id)
	}
	condition = fmt.Sprintf(condition, param...)
	query = fmt.Sprintf(query, model.ExpireAt, condition, model.StoreId)
	log.Debug(query)
	_, err := DB.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil

}
