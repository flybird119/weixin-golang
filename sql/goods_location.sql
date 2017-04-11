--create by lixiao
CREATE TABLE goods_location (
	id uuid primary key default gen_random_uuid(),
    goods_id text not null,
    type int not null,                                          --类型 0 旧 1 新
    storehouse_id text not null default '',                 --仓库id
    shelf_id text not null default '',                      --货架id
    floor_id text not null default '',                      --货架层id
	create_at timestamptz not null default now(),
	update_at timestamptz not null default now()
);
