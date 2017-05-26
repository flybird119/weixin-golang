create SEQUENCE store_renew_record_id;

CREATE TABLE store_renew_record(
	"id" text primary key not null default trim(to_char(nextval('store_renew_record_id'), '00000000')),
    charges int default 0,
    remark text default '',
    start_at timestamptz not null ,
    end_at timestamptz not null ,
    create_at timestamptz not null default now(),
	update_at timestamptz not null default now()
);
