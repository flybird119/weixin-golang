--create by lixiao

CREATE TABLE store_withdraw_card(
	id uuid primary key default gen_random_uuid(),      --代理主键
	store_id text not null,
	card_type int not null default 0,   				--0对私账户 1对公账户
	card_no text not null,								--银行卡账号
	card_name text not null, 							--银行卡所属银行
	username text not null, 							--用户名
	create_at timestamptz  not null default now(),
	update_at timestamptz  not null default now()
);
CREATE UNIQUE INDEX IF NOT EXISTS store_withdraw_card_store on store_withdraw_card(store_id);
