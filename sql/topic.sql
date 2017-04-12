create SEQUENCE if not exists topic_id_seq;

create table topic (
    id text primary key not null default 'TP' || to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('topic_id_seq'), '00000000')),

    profile text default '',
    title text not null,

    store_id text not null,
    sort int not null default 1,  -- 1 优先级低  2 优先级中   3 优先级高
    status int not null default 1, -- 1 正常  2 下架
    create_at timestamptz not null default now(),
    update_at timestamptz not null default now()
);
