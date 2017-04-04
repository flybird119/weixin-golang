\-- 实体店表
create table real_shop (
    id uuid primary key default gen_random_uuid(),

    -- 实体店名称
    name text not null,

    -- 实体店地址信息
    province_code int not null,
    city_code int not null,
    scope_code int not null default 0,
    address text not null,

    -- 店铺图片
    images text not null default '',

    -- 关联的云店铺ID
    store_id text not null,

    create_at timestamptz not null default now(),
    update_at timestamptz not null default now()
);
