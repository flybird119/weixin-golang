--create by lixiao
create SEQUENCE master_id_seq;

CREATE TABLE master(
	id text primary key not null default trim(to_char(nextval('master_id_seq'), '00000000')),
	mobile varchar(15) not null,
	password text not null,
	create_at timestamptz not null default now(),
	update_at timestamptz not null default now()
);
