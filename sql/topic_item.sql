create SEQUENCE if not exists topic_item_id_seq;

create table topic_item (
    id text primary key not null default trim(to_char(nextval('topic_item_id_seq'), '000000000')),
    topic_id text not null,
    goods_id text not null,
    status int not null default 1, -- 1 正常  2 下架
    create_at timestamptz not null default now()
);
