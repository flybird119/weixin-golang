-- 中国高校所有专业名称
CREATE SEQUENCE IF NOT EXISTS groupon_purchase_log_id_seq;
CREATE TABLE IF NOT EXISTS groupon_purchase_log (
  id                    TEXT      PRIMARY KEY       NOT NULL                 trim(to_char(nextval('groupon_purchase_log'), '00000000')),--id
  -- 关联的表结构
  groupon_id            TEXT                        NOT NULL                 ,          --班级购id
  user_id               TEXT                        NOT NULL                 DEFAULT '',--创建人姓名
  create_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now()
);
CREATE  INDEX IF NOT EXISTS groupon_purchase_log_groupon ON groupon_purchase_log (groupon_id);
CREATE  INDEX IF NOT EXISTS groupon_purchase_log_user ON groupon_purchase_log (user_id);
