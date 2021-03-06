create table account_item(
 id          UUID PRIMARY KEY         NOT NULL                 DEFAULT gen_random_uuid(),
 user_type          INT                 NOT NULL,                             --类型 1 商家店铺 2 平台
 store_id           TEXT                NOT NULL                 DEFAULT '',  --云店铺id
 order_id           TEXT                NOT NULL                 DEFAULT '',  --订单id
 remark             TEXT                NOT NULL                 DEFAULT '',  --备注
 --类型 商家 1待结算-交易完成 2 待结算-手续费 4 待结算-交易收入 17 可提现-交易完成 18 可提现-充值 20 可提现-体现 24 可提现-售后
 --类型 平台 33 平台-订单手续费
 item_type          INT                 NOT NULL,
 item_fee           INT                 NOT NULL                 DEFAULT 0,
 account_balance    BIGINT              NOT NULL,                             --当前账户余额
 create_at   TIMESTAMP WITH TIME ZONE NOT NULL                 DEFAULT now(),
 update_at   TIMESTAMP WITH TIME ZONE NOT NULL                 DEFAULT now()
);
CREATE INDEX IF NOT EXISTS account_item_store_id ON account_item (store_id);
