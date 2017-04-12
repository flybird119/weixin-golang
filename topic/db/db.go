package db

import (
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
		err = AddTopicItem(topic.Items[i])
		if err != nil {
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
	args = append(args, topic.Id)
	condition += fmt.Sprintf(" where id=$%d", len(args))

	args = append(args, topic.TokenStoreId)
	condition += fmt.Sprintf(" and store_id=$%d", len(args))
	query += condition
	_, err := DB.Exec(query, args)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//AddTopicItem 增加话题项 topicItem.TopicId topicItem.GoodsId
func AddTopicItem(topicItem *pb.TopicItem) error {
	query := "insert topic_item (topic_id,goods_id) values($1,$2)"
	log.Debugf("insert topic_item (topic_id,goods_id) values(%s,%s)", topicItem.TopicId, topicItem.GoodsId)
	_, err := DB.Exec(query, topicItem.TopicId, topicItem.GoodsId)
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
	_, err := DB.Exec(query, topicItem.Id)
	if err != nil {
		misc.LogErr(err)
		return err
	}
	return nil
}

//SearchTopics 搜索话题 topic.Id topic.Title topic.TokenStoreId
func SearchTopics(topic *pb.Topic) (topics []*pb.Topic, err error) {
	query := "select t.id,t.profile,t.title,t.sort,t.status, extract(epoch from t.create_at)::integer t.create_at,extract(epoch from t.update_at)::integer t.update_at,count(ti.id) from topic t join topic_item ti on ti.topic_id=t.id where 1=1"
	var args []interface{}
	var condition string
	if topic.Id != "" {
		args = append(args, topic.Id)
		condition += fmt.Sprintf(" and id=$%d", len(args))
	}
	if topic.Title != "" {
		args = append(args, topic.Title)
		condition += fmt.Sprintf(" and id=$%d", len(args))
	}
	if topic.TokenStoreId != "" {
		args = append(args, topic.TokenStoreId)
		condition += fmt.Sprintf(" and id=$%d", len(args))
	}

	condition += "order by t.sort"
	query += condition

	log.Debugf(query+" args:%s", args)
	rows, err := DB.Query(query, args)
	if err != nil {
		misc.LogErr(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		//select t.id,t.profile,t.title,t.sort,t.status, extract(epoch from t.create_at)::integer t.create_at,extract(epoch from t.update_at)::integer t.update_at,count(ti.id)
		topic := &pb.Topic{}
		topics = append(topics, topic)
		rows.Scan(&topic.Id, &topic.Profile, &topic.Title, &topic.Sort, &topic.Status, &topic.CreateAt, &topic.UpdateAt, &topic.ItemCount)
	}
	return topics, err
}
