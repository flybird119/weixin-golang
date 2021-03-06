-- 账户基本信息
CREATE TABLE IF NOT EXISTS account (
  id                 UUID PRIMARY KEY         NOT NULL                 DEFAULT gen_random_uuid(),

  type                  SMALLINT                 NOT NULL                         DEFAULT 1,       --1商家 2.平台
  balance               BIGINT                   NOT NULL                 DEFAULT 0,       --可提现余额
  unsettled_balance     BIGINT                   NOT NULL                 DEFAULT 0,
  store_id              text                     NOT NULL                 DEFAULT '',       --商家id ，平台就是0

  online_income         BIGINT                   NOT NULL                 DEFAULT 0,       --线上累计收入
  expenses              BIGINT                   NOT NULL                 DEFAULT 0,       --线上累计支出
  offline_income        BIGINT                   NOT NULL                 DEFAULT 0,       --线下累计收入

  create_at             TIMESTAMP WITH TIME ZONE NOT NULL                 DEFAULT now(),   --创建时间
  update_at             TIMESTAMP WITH TIME ZONE NOT NULL                 DEFAULT now()    --更新时间
);
CREATE  INDEX IF NOT EXISTS account_type ON account (type);
CREATE  INDEX IF NOT EXISTS account_store ON account (store_id);
CREATE UNIQUE INDEX IF NOT EXISTS account_id on account (id);
