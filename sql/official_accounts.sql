/* 记录已授权的微信公众号 */
create SEQUENCE official_accounts_id_seq;

create table official_accounts (
    "id" text primary key not null default trim(to_char(nextval('location_id_seq'), '00000000')),
    appid text not null,

    /* 授权方 昵称*/
    nick_name text not null,
    head_img text not null,
    user_name text not null,

    /* 授权方公众号的原始ID */
    principal_name text not null default '',
    qrcode_url text not null,
    service_type_info smallint not null,
    verify_type_info smallint not null,

    create_at timestamptz not null default now()
);
