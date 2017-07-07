-- 中国高校所有专业名称
CREATE SEQUENCE IF NOT EXISTS groupon_id_seq;
CREATE TABLE IF NOT EXISTS groupon (
  id                    TEXT      PRIMARY KEY       NOT NULL                 DEFAULT to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('groupon_id_seq'), '0000000')),--id
  status                INT                         NOT NULL                 DEFAULT 1, --1 正常使用 2 异常
  -- 关联的表结构
  store_id              TEXT                        NOT NULL                 ,          --店铺id
  school_id             TEXT                        NOT NULL                 ,          --学校id
  institute_id          TEXT                        NOT NULL                 ,          --学院id
  institute_major_id   TEXT                        NOT NULL                 ,          --专业id
  founder_id            TEXT                        NOT NULL                 ,          --专业id

  -- 团购的固有字段
  term                  TEXT                        NOT NULL                 DEFAULT '',--学期
  class                 TEXT                        NOT NULL                 DEFAULT '',--班级名称或者班级编号
  founder_type          SMALLINT                    NOT NULL ,                          -- 1 学生 2 商家
  founder_name          TEXT                        NOT NULL                 DEFAULT '',--创建人姓名
  founder_mobile        TEXT                        NOT NULL                 DEFAULT '',--创建人手机号
  profile               TEXT                        NOT NULL                 DEFAULT '',-- 简介
  participate_num       INT                         NOT NULL                 DEFAULT 0, --参与人数
  star_num              INT                         NOT NULL                 DEFAULT 0, --点赞的人数
  total_sales           BIGINT                      NOT NULL                 DEFAULT 0,   --总销售额
  order_num             INT                         NOT NULL                 DEFAULT 0, --订单数量
  create_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now(),
  expire_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 ,          --失效时间
  update_at             TIMESTAMP WITH TIME ZONE    NOT NULL                 DEFAULT now()
);
