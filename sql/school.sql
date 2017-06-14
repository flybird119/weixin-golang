create table school (
    id uuid primary key default gen_random_uuid(),

    -- 基本信息数据
    name text not null,
    tel text not null,
    express_fee int not null,
    store_id text not null,

    -- 位置信息
    lat double precision NOT NULL DEFAULT 0,
    lng double precision NOT NULL DEFAULT 0,
    status int default 0, --学校状态 0 正常 1异常
    del_at timestamptz,
    del_staff_id text default '',
    is_recyling bool default true,
    create_at timestamptz not null default now(),
	update_at timestamptz not null default now()
);
