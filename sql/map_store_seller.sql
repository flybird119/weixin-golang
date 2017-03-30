--create by lixiao 
CREATE TABLE map_store_seller (
	id serial,
	role int8,
	store_id int4,
	seller_id int4,
	status int4,
	PRIMARY KEY (id)
)
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."id" IS '代理主键';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."role" IS '权限';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."store_id" IS '店铺id';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."seller_id" IS '商家id';
-- COMMENT ON COLUMN "public"."bc_map_store_seller"."status" IS '状态';
