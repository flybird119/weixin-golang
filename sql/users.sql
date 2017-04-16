CREATE SEQUENCE if not exists users_id;
create table users (
    id text primary key not null default trim(to_char(nextval('users_id'), '0000000000')),
    openid text not null,
    nickname text not null,
    sex int not null default 3,               --1男 2女 其他：未知
    avatar text not null,
    status int default 1,           --1正常 2 异常
    create_at timestamptz not null default now(),
    update_at timestamptz not null default now()
);
