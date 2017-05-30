-- 订单表
CREATE SEQUENCE IF NOT EXISTS withdrawals_id_seq;

create table withdrawals (
    id text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('withdrawals_id_seq'), '00000000')),
    -- 店铺信息
    store_id text not null,                            --店铺id
    withdraw_card_id text not null,                    --学校id
    staff_id text not null,
    --银行卡信息
    card_type int not null default 0,                  --0 对私账户 1 对公账户
    card_no text not null,                             --银行卡账号
    card_name text not null,                           --银行卡所属银行
    username text not null,                         --用户名
    apply_phone text default '',
    -- 金额数据
    withdraw_fee int not null,                         --提现数据
    --提现状态
    status int not null default 1,                     --提现状态 1待处理 2处理中 3处理成功
    --操作时间
    admin_id text default '',
    apply_at timestamptz not null default now(),       --创建时间
    complete_at timestamptz,
    accept_at timestamptz,                            --完成时间
    update_at timestamptz not null default now()       --更新时间

);
