
create table statistic_goods_sales(
    id          UUID PRIMARY KEY                            NOT NULL                 DEFAULT gen_random_uuid(), --代理主键
    --店铺 学校 id
    store_id                        TEXT                    NOT NULL,                           --云店铺id
    school_id                       TEXT                    NOT NULL,                           --学校id

    --线上统计
    alipay_order_num                INT                     NOT NULL            DEFAULT 0,      --当天支付宝订单量
    alipay_order_fee                INT                     NOT NULL            DEFAULT 0,      --当天支付宝订单销售额
    wechat_order_num                INT                     NOT NULL            DEFAULT 0,      --当天微信订单销售量
    wechat_order_fee                INT                     NOT NULL            DEFAULT 0,      --当天微信订单销售额
    online_new_book_sales_fee       INT                     NOT NULL            DEFAULT 0,      --当天线上新书销售额
    online_old_book_sales_fee       INT                     NOT NULL            DEFAULT 0,      --当天线上旧书销售额
    send_order_num                  INT                     NOT NULL            DEFAULT 0,      --当天配送的订单量（订单已发货）

    --售后统计
    after_sale_num                  INT                     NOT NULL            DEFAULT 0,      --当天售后新增数量
    after_sale_handled_num          INT                     NOT NULL            DEFAULT 0,      --售后处理数量
    after_sale_handled_fee          INT                     NOT NULL            DEFAULT 0,      --售后处理金额

    --线下统计
    offline_new_book_sales_fee      INT                     NOT NULL            DEFAULT 0,      --线下新书销售额
    offline_old_book_sales_fee      INT                     NOT NULL            DEFAULT 0,      --线下旧书销售额
    offline_order_num               INT                     NOT NULL            DEFAULT 0,      --线下订单数量

    --时间字段
    statistic_at                    TIMESTAMP WITH TIME ZONE NOT NULL  ,
    create_at                       TIMESTAMP WITH TIME ZONE NOT NULL           DEFAULT now()  --创建时间                      --改数据统计所属年
);
CREATE INDEX IF NOT EXISTS statistic_goods_school ON statistic_goods_sales(school_id);
CREATE INDEX IF NOT EXISTS statistic_goods_store ON statistic_goods_sales(store_id);
CREATE UNIQUE INDEX IF NOT EXISTS statistic_goods_id on statistic_goods_sales(id);
