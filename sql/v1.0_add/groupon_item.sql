-- 中国高校所有专业名称
CREATE SEQUENCE IF NOT EXISTS groupon_item_id_seq;
CREATE TABLE IF NOT EXISTS groupon_item (
  id                    TEXT      PRIMARY KEY       NOT NULL                 DEFAULT trim(to_char(nextval('groupon_item_id_seq'), '00000000')),--id
  -- 关联的表结构
  groupon_id            TEXT                        NOT NULL                 ,          --班级购id
  goods_id              TEXT                        NOT NULL                 ,          --商品i
  create_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now(),
  update_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now()
);
