-- 订单表
CREATE SEQUENCE IF NOT EXISTS recharge_id_seq;
create table recharge (
    id text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('recharge_id_seq'), '0000000')),
    -- 店铺信息
    store_id text not null,                            --店铺id
    --充值金额
    recharge_fee int not null,                         --充值金额
    pay_way text not null,
    trade_no text,
    --提现状态
    status int not null default 1,                     --提现状态 1申请中 2 处理成功
    --操作时间
    apply_at timestamptz not null default now(),       --创建时间
    complete_at timestamptz,                           --完成时间
    update_at timestamptz not null default now()       --更新时间
);
