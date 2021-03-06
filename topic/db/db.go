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

//AddTopic 增加话题 topic.Profile topic.Title topic.TokenStoreId topic.Sort topic.Items
func AddTopic(topic *pb.Topic) error {
	//首先保存话题，然后保存话题项
	query := "insert into topic (profile,title,store_id,sort) values($1,$2,$3,$4) returning id"
	log.Debugf("insert into topic (profile,title,store_id,sort) values(%s,%s,%s,%d) returning id", topic.Profile, topic.Title, topic.TokenStoreId, topic.Sort)
	err := DB.QueryRow(query, topic.Profile, topic.Title, topic.TokenStoreId, topic.Sort).Scan(&topic.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	errStr := ""
	for i := 0; i < len(topic.Items); i++ {
		topic.Items[i].TopicId = topic.Id
		err = AddTopicItem(topic.Items[i])
		if err != nil {
			misc.LogErr(err)
			errStr += fmt.Sprintf("第%d项上传时失败", i)
		}
	}
	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}

//DelTopic 删除话题 topic.Id topic.TokenStoreId
func DelTopic(topic *pb.Topic) error {
	tx, err := DB.Begin()
	if err != nil {
		misc.LogErr(err)
	}
	defer tx.Rollback()
	//首先删除topic
	query := "delete from topic where id=$1 and store_id=$2"
	log.Debugf("delete from topic where id=%s and store_id=%s", topic.Id, topic.TokenStoreId)
	_, err = tx.Exec(query, topic.Id, topic.TokenStoreId)
	if err != nil {
		misc.LogErr(err)
		return err
	}

	//删除topic item
	query = "delete from topic_item where topic_id=$1"
	log.Debugf("delete from topic_item where topic_id=%s", topic.Id)
	_, err = tx.Exec(query, topic.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	tx.Commit()
	//再删除
	return nil
}

//UpdateTopic 更新话题 topic.Profile topic.Title topic.Sort topic.Id
func UpdateTopic(topic *pb.Topic) error {
	query := "update topic set update_at=now()"
	var args []interface{}
	var condition string
	if topic.Profile != "" {
		args = append(args, topic.Profile)
		condition += fmt.Sprintf(",profile=$%d", len(args))
	}
	if topic.Title != "" {
		args = append(args, topic.Title)
		condition += fmt.Sprintf(",title=$%d", len(args))
	}
	if topic.Sort != 0 {
		args = append(args, topic.Sort)
		condition += fmt.Sprintf(",sort=$%d", len(args))
	}
	if topic.Status != 0 {
		args = append(args, topic.Status)
		condition += fmt.Sprintf(",status=$%d", len(args))
	}

	args = append(args, topic.Id)
	condition += fmt.Sprintf(" where id=$%d", len(args))

	args = append(args, topic.TokenStoreId)
	condition += fmt.Sprintf(" and store_id=$%d", len(args))
	query += condition
	_, err := DB.Exec(query, args...)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//AddTopicItem 增加话题项 topicItem.TopicId topicItem.GoodsId
func AddTopicItem(topicItem *pb.TopicItem) error {
	query := "select goods_id,count(goods_id) from topic_item where topic_id=$1 group by goods_id"
	rows, err := DB.Query(query, topicItem.TopicId)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var goods_id string
		var count int
		err = rows.Scan(&goods_id, &count)
		if err != nil {
			misc.LogErr(err)
			return err
		}
		if goods_id == topicItem.GoodsId {

			return errors.New("不能更添加重复的商品")
		}

	}
	query = "insert into topic_item (topic_id,goods_id) values($1,$2)"
	log.Debugf("insert into topic_item (topic_id,goods_id) values(%s,%s)", topicItem.TopicId, topicItem.GoodsId)
	_, err = DB.Exec(query, topicItem.TopicId, topicItem.GoodsId)
	if err != nil {
		log.Errorf("%+v", err)
		return err
	}
	return nil

}

//DelTopicItem 删除话题项 topicItem.Id
func DelTopicItem(topicItem *pb.TopicItem) error {
	query := "delete from topic_item where id=$1 and topic_id=$2"
	log.Debugf("delete from topic_item where id=%s and topic_id=%s", topicItem.Id, topicItem.TopicId)
	_, err := DB.Exec(query, topicItem.Id, topicItem.TopicId)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//SearchTopics 搜索话题 topic.Id topic.Title topic.TokenStoreId
func SearchTopics(topic *pb.Topic, searchType int64) (topics []*pb.Topic, err error, totalCount int64) {
	query := "select t.id,t.profile,t.title,t.sort,t.status, extract(epoch from t.create_at)::bigint create_at,extract(epoch from t.update_at)::bigint update_at  from topic t where exists (select * from topic_item ti where ti.topic_id=t.id) "
	countQuery := "select count(*) from topic t where exists (select * from topic_item ti where ti.topic_id=t.id) "
	var args []interface{}
	var condition string
	if topic.Id != "" {
		args = append(args, topic.Id)
		condition += fmt.Sprintf(" and t.id=$%d", len(args))
	}
	if topic.Title != "" {
		args = append(args, misc.FazzyQuery(topic.Title))
		condition += fmt.Sprintf(" and t.title like $%d", len(args))
	}
	if topic.TokenStoreId != "" {
		args = append(args, topic.TokenStoreId)
		condition += fmt.Sprintf(" and t.store_id=$%d", len(args))
	}
	if topic.Status != 0 {
		args = append(args, topic.Status)
		condition += fmt.Sprintf(" and t.status=$%d", len(args))
	}
	var page, size, topic_page, topic_size int64

	if searchType == 1 {
		page = 1
		size = 100
		topic_page = 1
		topic_size = 100
	} else {
		if topic.Page <= 0 {
			page = 1
		} else {
			page = topic.Page
		}
		if topic.Size <= 0 {
			size = 15
		} else {
			size = topic.Size
		}
		if topic.TopicPage <= 0 {
			topic_page = 1
		} else {
			topic_page = topic.TopicPage
		}
		if topic.TopicSize <= 0 {
			topic_size = 15
		} else {
			topic_size = topic.TopicSize
		}
	}

	//计数count
	countQuery += condition
	err = DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return
	}
	//如果统计的为零
	if totalCount == 0 {
		return
	}

	condition += fmt.Sprintf(" order by t.sort desc,t.id desc  OFFSET %d LIMIT %d", (topic_page-1)*topic_size, topic_size)
	query += condition

	log.Debugf(query+" args:%s", args)
	rows, err := DB.Query(query, args...)
	if err == sql.ErrNoRows {
		return nil, nil, totalCount
	}
	if err != nil {
		misc.LogErr(err)
		return nil, err, totalCount
	}
	defer rows.Close()
	for rows.Next() {
		//select t.id,t.profile,t.title,t.sort,t.status, extract(epoch from t.create_at)::bigint t.create_at,extract(epoch from t.update_at)::bigint t.update_at,count(ti.id)
		topic := &pb.Topic{}
		topics = append(topics, topic)
		rows.Scan(&topic.Id, &topic.Profile, &topic.Title, &topic.Sort, &topic.Status, &topic.CreateAt, &topic.UpdateAt)

		items, findErr, totalCount1 := GetTopicItemsByTopic(topic.Id, page, size)
		if err != nil {

			misc.LogErr(findErr)
			return nil, findErr, totalCount
		}

		topic.ItemCount = int64(totalCount1)

		topic.Items = items
	}
	return topics, err, totalCount
}

//GetTopicItemsByTopic 获取话题项
func GetTopicItemsByTopic(topic_id string, page, size int64) (items []*pb.TopicItem, err error, totalCount int64) {
	query := fmt.Sprintf("select count(*) from topic_item t where topic_id='%s' and exists(select * from goods g where g.id=t.goods_id and (g.has_new_book=true or g.has_old_book=true))", topic_id)
	log.Debug(query)
	err = DB.QueryRow(query).Scan(&totalCount)
	if err != nil {
		misc.LogErr(err)
		return
	}
	query = "select t.id,t.topic_id,t.goods_id,t.status,extract(epoch from t.create_at)::bigint create_at from topic_item t where topic_id=$1 and exists(select * from goods g where g.id=t.goods_id and (g.has_new_book=true or g.has_old_book=true)) order by id"
	log.Debugf("select id,topic_id,goods_id,status,extract(epoch from create_at)::bigint create_at from topic_item t where topic_id='%s' order by id OFFSET %d LIMIT %d ", topic_id, (page-1)*size, size)
	query += fmt.Sprintf(" OFFSET %d LIMIT %d ", (page-1)*size, size)

	rows, err := DB.Query(query, topic_id)

	if err == sql.ErrNoRows {
		return nil, nil, 0
	}
	if err != nil {
		misc.LogErr(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &pb.TopicItem{}
		items = append(items, item)
		err = rows.Scan(&item.Id, &item.TopicId, &item.GoodsId, &item.Status, &item.CreateAt)

		if err != nil {
			misc.LogErr(err)
			return
		}
	}
	return
}
