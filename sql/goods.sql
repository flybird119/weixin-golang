create sequence goods_id_seq;

create table goods(

    id text primary key not null default trim(to_char(nextval('goods_id_seq'),'0000000000')),
    book_id text not null,                                  --bookid
    store_id text not null,                                 --云店铺id
    isbn text not null,                                     --isbn
    new_book_amount int not null default 0,                 --新书数量
    old_book_amount int not null default 0,                 --旧书数量
    new_book_price  int not null default 0,                 --新书价格
    old_book_price int not null default 0,                  --旧书价格
    has_new_book bool default false,                        --是否有新书
    has_old_book bool default false,                        --是否有旧书
    create_at timestamptz not null default now(),           --创建时间
    update_at timestamptz not null default now(),           --更改时间
    new_book_sale_amount int not null default 0,            --新书销售数量
    old_book_sale_amount int not null default 0,            --二手书销售数量
    is_selling bool default true                            --在售状态
)
