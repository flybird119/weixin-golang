-- 中国高校所有专业名称
CREATE SEQUENCE IF NOT EXISTS shared_major_id_seq;
CREATE TABLE IF NOT EXISTS shared_major (
  id                    TEXT      PRIMARY KEY       NOT NULL                 DEFAULT trim(to_char(nextval('shared_major_id_seq'),'000000000')),--id
  no                    TEXT                        NOT NULL                 DEFAULT '',
  name                  TEXT                        NOT NULL                 DEFAULT '',
  create_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now(),
  update_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now()
);
CREATE  INDEX IF NOT EXISTS shared_major_name ON shared_major (name);
