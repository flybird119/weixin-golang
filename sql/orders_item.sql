-- 订单项
create table orders_item (
    id uuid primary key default gen_random_uuid(),

    -- 关联项
    goods_id text not null,     --商品id
    orders_id text,             --订单id
    type int not null,          --图书类型
    amount int not null,        --购买数量
    price int not null,         --商品单价
    create_at timestamptz not null default now() --创建时间
);

CREATE INDEX order_item_order ON orders_item(orders_id);
CREATE INDEX order_item_goods ON orders_item(goods_id);
