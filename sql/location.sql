-- 记录图书位置

create SEQUENCE location_id_seq;

create table location (
    "id" text primary key not null default trim(to_char(nextval('location_id_seq'), '00000000')),

    -- 辅助查找信息
    level smallint not null default 0,
    pid text not null default '',
    store_id text not null,
    -- 地址名称
    name text not null,

    -- 创建时间、修改时间
    create_at timestamptz not null default now(),
	update_at timestamptz not null default now()
);
