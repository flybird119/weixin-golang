create table statistic_book_sales(
    id          UUID PRIMARY KEY                            NOT NULL                 DEFAULT gen_random_uuid(), --代理主键
    --店铺 学校 id
    store_id                        TEXT                    NOT NULL,                               --云店铺id
    goods_id                        TEXT                    NOT NULL,                               --学校id

    online_new_book_sales_num       INT                     NOT NULL                DEFAULT 0,      --小伤新书销售量
    online_old_book_sales_num       INT                     NOT NULL                DEFAULT 0,      --线上旧书销售量
    offline_new_book_sales_num      INT                     NOT NULL                DEFAULT 0,      --线下新书销售量
    offline_old_book_sales_num      INT                     NOT NULL                DEFAULT 0,      --线下旧书销售量

    new_book_remove_num             INT                     NOT NULL                DEFAULT 0,      --新书下架量
    new_book_upload_num             INT                     NOT NULL                DEFAULT 0,      --新书上架量
    old_book_remove_num             INT                     NOT NULL                DEFAULT 0,      --旧书下架量
    old_book_upload_num             INT                     NOT NULL                DEFAULT 0,      --旧书上架量
    --时间字段
    create_at                       TIMESTAMP WITH TIME ZONE NOT NULL               DEFAULT now(),  --创建时间
    statistic_day                   TIMESTAMP WITH TIME ZONE NOT NULL,                              --改数据统计所属天
    statistic_month                 TIMESTAMP WITH TIME ZONE NOT NULL,                              --改数据统计所属月份
    statistic_year                  TIMESTAMP WITH TIME ZONE NOT NULL                               --改数据统计所属年

)
