-- 订单表
CREATE SEQUENCE IF NOT EXISTS retail_id_seq;

create table retail (
    id text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('retail_id_seq'), '00000000')),
    -- 金额数据
    total_fee int not null,                            --用户支付费用
    goods_fee int not null,                            --商品费用
    -- 店铺信息
    store_id text not null,                            --店铺id
    school_id text not null,                           --学校id
    --操作人
    handle_staff_id text default '',                   --处理员工id
    create_at timestamptz not null default now(),                             --创建时间
    update_at timestamptz not null default now()       --更新时间
);

CREATE UNIQUE INDEX retail_id on retail(id); --建立id索引
CREATE INDEX retail_store on retail(store_id); --建立id索引
CREATE INDEX retail_school on retail(school_id); --建立id索引
