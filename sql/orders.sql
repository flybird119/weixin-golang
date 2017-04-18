-- 订单表
create SEQUENCE order_id_seq;

create table orders (
    "id" text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('order_id_seq'), '00000000')),
    order_status smallint default 0,    -- 订单状态

    -- 金额数据
    total_fee int not null,             --用户支付费用
    freight int default 0,              --运费
    goods_fee int not null,             --商品费用
    withdrawal_amount int not null      --可体现金额

    -- 关联信息
    user_id text not null,              --用户id
    mobile text not null,               --联系人手机号
    name text not null,                 --联系人手机号
    address text not null,              --联系人地址
    remark text default '',             --订单备注

    -- 店铺信息
    store_id text not null,             --店铺id
    school_id text not null,            --学校id

    -- 支付宝或微信信息
    trade_no text default '',           --第三方交易号

    --订单时间信息
    order_at timestamptz not null default now(),    --下单时间
    pay_at timestamptz,                             --付款时间
    deliver_at timestamptz,                         --发货时间
    print_at timestamptz,                           --打印时间
    complete_at timestamptz,                        --完成时间

    --操作人
    print_staff_id text,                            --打印员工id
    deliver_staff_id text,                          --发货负责人
    after_sale_staff_id text,                       --售后负责人

    --售后
    after_sale_apply_at timestamptz,                --售后开始时间
    after_sale_end_at timestamptz,                  --售后结束时间
    after_sale_status int default 0,                --售后单号、状态 1待处理 2 退款中 3退款失败 4退款成功
    after_sale_trad_no text,                        --售后交易号
    refund_fee int default 0,                       --退款金额


    update_at timestamptz not null default now()    --更新时间

);
