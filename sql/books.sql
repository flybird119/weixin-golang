-- 存储标准图书信息
create SEQUENCE book_id_seq;

create table books (
    "id" text primary key not null default trim(to_char(nextval('book_id_seq'), '00000000')),

    /* 附加信息记录 */
    store_id text not null,
    level smallint not null default 0,

    /* 图书主信息 */
    title text not null,
    isbn text not null,
    price int not null,
    author text not null default '',
    publisher text not null default '',
    pubdate text not null default '',
    subtitle text not null default '',
    image text not null default '',
    summary text not null default '',
    author_intro text not null default '',
    /* 时间记录 */
    create_at timestamptz not null default now()
);
