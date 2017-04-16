CREATE SEQUENCE if not exists cart_id_seq;
--购书车数据库表
create table cart(
id uuid primary key default gen_random_uuid(),      --代理主键
user_id text not null,                              --用户id
store_id text not null,                             --云店铺id
goods_id text not null,                             --商品id
type int not null,                                  --0 新书 1 旧书
amount int not null,                                --购买数量
create_at timestamptz not null default now()        --创建时间
)
