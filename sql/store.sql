--create by lixiao
create SEQUENCE store_id_seq;

CREATE TABLE store(
	"id" text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('store_id_seq'), '000000')),
	name varchar(50) not null,
	logo text ,
	status int default 0,
	profile text,
	create_at timestamptz  not null default now(),
	service_mobiles text ,
	expire_at timestamptz  not null,
	address text ,
	map_address text ,
	business_license text
)
-- COMMENT ON COLUMN "public"."bc_info_store"."name" IS '店铺名称';
-- COMMENT ON COLUMN "public"."bc_info_store"."logo" IS '店铺logo';
-- COMMENT ON COLUMN "public"."bc_info_store"."status" IS '店铺状态';
-- COMMENT ON COLUMN "public"."bc_info_store"."profile" IS '店铺简介';
-- COMMENT ON COLUMN "public"."bc_info_store"."create_time" IS '创建时间';
-- COMMENT ON COLUMN "public"."bc_info_store"."service_mobiles" IS '服务电话';
-- COMMENT ON COLUMN "public"."bc_info_store"."expire_time" IS '失效时间';
-- COMMENT ON COLUMN "public"."bc_info_store"."address" IS '店铺地址';
-- COMMENT ON COLUMN "public"."bc_info_store"."map_address" IS '店铺地图地址';
-- COMMENT ON COLUMN "public"."bc_info_store"."business_license" IS '营业执照';
