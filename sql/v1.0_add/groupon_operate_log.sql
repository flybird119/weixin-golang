-- 中国高校所有专业名称
CREATE SEQUENCE IF NOT EXISTS groupon_operate_log_id_seq;
CREATE TABLE IF NOT EXISTS groupon_operate_log (
  id                    TEXT      PRIMARY KEY       NOT NULL                 trim(to_char(nextval('groupon_operate_log_id_seq'), '00000000')),--id
  -- 关联的表结构
  groupon_id            TEXT                        NOT NULL                 ,          --班级购id
  founder_id            TEXT                        NOT NULL                 ,          --创建人id
  founder_type          SMALLINT                    NOT NULL ,                          -- 1 学生 2 商家
  founder_name          TEXT                        NOT NULL                 DEFAULT '',--创建人姓名
  operate_type          TEXT                        NOT NULL                 DEFAULT '',--操作类型  create  update   share  purchase
  operate_detail        TEXT                        NOT NULL                 DEFAULT '',--操作具体详情
  create_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now(),
  update_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now()
);
