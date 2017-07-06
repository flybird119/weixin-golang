-- 中国高校所有专业名称
CREATE SEQUENCE IF NOT EXISTS map_school_institute_id_seq;
CREATE TABLE IF NOT EXISTS map_school_institute (
  id                    TEXT      PRIMARY KEY    NOT NULL                 DEFAULT trim(to_char(nextval('map_school_institute_id_seq'),'000000000')),--id
  school_id             TEXT                     NOT NULL                 DEFAULT '',
  name                  TEXT                     NOT NULL                 DEFAULT '',
  create_at             TIMESTAMP WITH TIME ZONE NOT NULL                 DEFAULT now(),
  update_at             TIMESTAMP WITH TIME ZONE NOT NULL                 DEFAULT now()
);
CREATE  INDEX IF NOT EXISTS map_school_institute_name ON map_school_institute (name);
