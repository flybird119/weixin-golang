-- 地址管理表,存储用户地址

create table address (
    id uuid primary key default gen_random_uuid(),
    name varchar(50) not null,
    tel varchar(20) not null,
    address text not null,

    -- 用户ID
    user_id text not null,

    --是否是默认状态
    is_default bool not null default false,
    create_at timestamptz not null default now(),
    update_at timestamptz not null default now()

);
CREATE UNIQUE INDEX user_address_id on address(id); --建立id索引
CREATE INDEX user_address_user on address(user_id); --建立用户索引
