--create by lixiao
CREATE TABLE map_store_seller (
	id uuid primary key default gen_random_uuid(),
	role int not null,
	store_id text not null default '',
	seller_id text not null default '',
	create_at timestamptz not null default now(),
	update_at timestamptz not null default now(),
	PRIMARY KEY (id)
)
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."id" IS '代理主键';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."role" IS '权限';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."store_id" IS '店铺id';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."seller_id" IS '商家id';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."status" IS '状态';
