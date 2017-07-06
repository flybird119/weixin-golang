-- 中国高校所有专业名称
create sequence IF NOT EXISTS map_institute_major_id;
CREATE TABLE IF NOT EXISTS map_institute_major (
  id                    TEXT      PRIMARY KEY       NOT NULL                 DEFAULT trim(to_char(nextval('map_institute_major_id'),'000000000')),--id
  institute_id          TEXT                        NOT NULL                 DEFAULT '',
  name                  TEXT                        NOT NULL                 DEFAULT '',
  create_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now(),
  update_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now()
);
CREATE  INDEX IF NOT EXISTS map_institute_major_name ON map_institute_major (name);
CREATE  INDEX IF NOT EXISTS map_institute_major_id ON map_institute_major (institute_id);
