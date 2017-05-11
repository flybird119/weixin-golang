/* 连接店铺、用户信息的中间表 */

create table map_store_users (
    store_id text not null,
    user_id text not null,
    openid text not null
);
