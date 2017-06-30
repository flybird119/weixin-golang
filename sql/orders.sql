-- 订单表
CREATE SEQUENCE IF NOT EXISTS order_id_seq;

create table orders (
    id text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('order_id_seq'), '00000000')),
    order_status smallint default 0,    -- 订单状态

    -- 金额数据
    total_fee int not null,             --用户支付费用
    freight int default 0,              --运费
    goods_fee int not null,             --商品费用
    withdrawal_fee int not null ,    --可体现金额

    -- 关联信息
    user_id text not null default '',              --用户id
    mobile text not null  default '',               --联系人手机号
    name text not null default '',                 --联系人手机号
    address text not null,              --联系人地址
    remark text default '',             --订单备注
    seller_remark text not null default '',      --商家备注 --新增模块
    seller_remark_type int not null default 0,  --商家备注类型
    -- 店铺信息
    store_id text not null,             --店铺id
    school_id text not null,            --学校id

    -- 支付宝或微信信息
    trade_no text default '',           --第三方交易号
    pay_channel text default '',        --alipay or wechat

    --订单时间信息
    order_at timestamptz not null default now(),    --下单时间
    pay_at timestamptz,                             --付款时间
    deliver_at timestamptz,                         --发货时间
    print_at timestamptz,                           --打印时间
    complete_at timestamptz ,                        --完成时间
    confirm_at timestamptz,                         --订单确认时间
    distribute_at timestamptz,                      --配送时间
    close_at timestamptz,                           --订单关闭时间
    after_sale_apply_at timestamptz,                --售后开始时间
    after_sale_end_at timestamptz,                  --售后结束时间

    --操作人
    print_staff_id text default '',                            --打印员工id
    deliver_staff_id text default '',                          --发货负责人
    distribute_staff_id text default '',                       --配送人员
    after_sale_staff_id text default '',                       --售后负责人

    --售后
    after_sale_status int default 0,                --售后单号、状态 1待处理 2 退款中 3退款失败 4退款成功
    after_sale_trad_no text default '',                        --售后交易号
    refund_fee int default 0,                       --退款金额
    apply_refund_fee int default 0,                 --申请退款金额
    after_sale_reason text default '',
    after_sale_images jsonb default '[]',

    --团购
    groupon_id text  default '',                                --班级购id

    update_at timestamptz not null default now()    --更新时间

);

CREATE  INDEX orders_store ON orders(store_id);
CREATE  INDEX orders_school ON orders(school_id);

-- >订单的状态由二进制来管理
-- 解释：<br>
--		 	 8  7  6  5  4  3  2  1 <br>
--	        _  _  _  _  _  _  _  _ <br>
--> 第一位: 是否支付 0 未支付 1 已支付 <br>
--> 第二位: 是否发货 0 未发货 1 已发货 <br>
--> 第三位: 是否完成 0 未完成 1 已完成 <br>
--> 第四位: 是否关闭 0 未关闭 1 已关闭 <br>
--> 第五位: 是否售后 0 未售后 1 售后状态 <br>
--> 未支付订单 0
--> 代发货订单 1
--> 已发货订单 3
--> 已完成订单 7
--> 已关闭订单 8
--> 售后订单 17 - 23 :<br>
-->>当前订单的状态n ,查看售后的订单进行到哪一步 n-16 匹配上面的值:<br>
