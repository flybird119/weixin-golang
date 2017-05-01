-- 订单项
create table retail_item (
    id uuid primary key default gen_random_uuid(),
    -- 关联项
    goods_id text not null,     --商品id
    retail_id text not null,    --订单id
    type int not null,          --图书类型
    amount int not null,        --购买数量
    price int not null,         --商品单价
    create_at timestamptz not null default now() --创建时间
);
CREATE UNIQUE INDEX retail_item_id on retail_item(id); --建立id索引
CREATE INDEX retail_item_retail ON retail_item(retail_id);
CREATE INDEX retail_item_goods ON retail_item(goods_id);
