--create by lixiao
CREATE TABLE seller(
	id serial,
	mobile varchar(15),
	password text,
	username varchar(24),
	name varchar(50),
	avatar text,
	create_at int,
	update_at int,
	status int,
	id_card varchar(50),
	PRIMARY KEY (id)
)
-- COMMENT ON TABLE "public"."bc_info_seller" IS '购书云商家注册表';
-- COMMENT ON COLUMN "public"."bc_info_seller"."id" IS '代理主键';
-- COMMENT ON COLUMN "public"."bc_info_seller"."mobile" IS '注册手机号';
-- COMMENT ON COLUMN "public"."bc_info_seller"."password" IS '登录密码';
-- COMMENT ON COLUMN "public"."bc_info_seller"."username" IS '用户昵称';
-- COMMENT ON COLUMN "public"."bc_info_seller"."name" IS '真实姓名';
-- COMMENT ON COLUMN "public"."bc_info_seller"."avatar" IS '头像';
-- COMMENT ON COLUMN "public"."bc_info_seller"."register_time" IS '注册时间';
-- COMMENT ON COLUMN "public"."bc_info_seller"."status" IS '状态';
-- COMMENT ON COLUMN "public"."bc_info_seller"."id_card" IS '身份证';
