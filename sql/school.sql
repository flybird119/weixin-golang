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

    create_at timestamptz not null default now(),
	update_at timestamptz not null default now()
);
