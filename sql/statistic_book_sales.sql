create table statistic_book_sales(
    id          UUID PRIMARY KEY                            NOT NULL                 DEFAULT gen_random_uuid(), --代理主键
    --店铺 学校 id
    store_id                        TEXT                    NOT NULL,                               --云店铺id
    school_id                       TEXT                    NOT NULL,                           --学校id
    goods_id                        TEXT                    NOT NULL,                               --学校id

    online_new_book_sales_num       INT                     NOT NULL                DEFAULT 0,      --小伤新书销售量
    online_old_book_sales_num       INT                     NOT NULL                DEFAULT 0,      --线上旧书销售量
    offline_new_book_sales_num      INT                     NOT NULL                DEFAULT 0,      --线下新书销售量
    offline_old_book_sales_num      INT                     NOT NULL                DEFAULT 0,      --线下旧书销售量
    --时间字段
    statistic_at                    TIMESTAMP WITH TIME ZONE NOT NULL  ,
    create_at                       TIMESTAMP WITH TIME ZONE NOT NULL               DEFAULT now()  --创建时间
);

CREATE UNIQUE INDEX IF NOT EXISTS statistic_books_store ON statistic_book_sales(store_id);
CREATE UNIQUE INDEX IF NOT EXISTS statistic_books_goods ON statistic_book_sales(goods_id);
