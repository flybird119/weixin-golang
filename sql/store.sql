--create by lixiao
CREATE TABLE store(
	id serial,
	name varchar(50) not null,
	logo text not null,
	stutus int4 default 0,
	profile text,
	create_time int4 not null,
	service_mobiles text not null,
	expire_time int4 not null,
	address text not null,
	map_address text not null,
	business_license text,
	PRIMARY KEY (id)
)
-- COMMENT ON COLUMN "public"."bc_info_store"."name" IS '店铺名称';
-- COMMENT ON COLUMN "public"."bc_info_store"."logo" IS '店铺logo';
-- COMMENT ON COLUMN "public"."bc_info_store"."stutus" IS '店铺状态';
-- COMMENT ON COLUMN "public"."bc_info_store"."profile" IS '店铺简介';
-- COMMENT ON COLUMN "public"."bc_info_store"."create_time" IS '创建时间';
-- COMMENT ON COLUMN "public"."bc_info_store"."service_mobiles" IS '服务电话';
-- COMMENT ON COLUMN "public"."bc_info_store"."expire_time" IS '失效时间';
-- COMMENT ON COLUMN "public"."bc_info_store"."address" IS '店铺地址';
-- COMMENT ON COLUMN "public"."bc_info_store"."map_address" IS '店铺地图地址';
-- COMMENT ON COLUMN "public"."bc_info_store"."business_license" IS '营业执照';
