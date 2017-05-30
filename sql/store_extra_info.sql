--create by lixiao
create SEQUENCE store_extra_seq;

CREATE TABLE store_extra_info(
	id text primary key not null default to_char(now() AT TIME ZONE 'cct', 'yymmdd') || trim(to_char(nextval('store_extra_seq'), '000000')),
	store_id text not null,
	poundage int default 0,
	charges int default 0,
	intention int default 1, --1未确认意愿 2有购买意愿 3无购买意愿 4已购买
	remark text default '',
	create_at timestamptz  not null default now(),
	update_at timestamptz not null default now()
);

CREATE UNIQUE INDEX IF NOT EXISTS store_extra_id ON store_extra_info(id);
CREATE UNIQUE INDEX IF NOT EXISTS store_extra_store ON store_extra_info(store_id);
