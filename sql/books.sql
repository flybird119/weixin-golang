-- 存储标准图书信息
create table books (
    id uuid primary key default gen_random_uuid(),

    /* 附加信息记录 */
    store_id text not null,
    level smallint not null default 0,

    /* 图书主信息 */



    /* 时间记录 */
    create_at timestamptz not null default now(),
    update_at timestamptz not null default now()
);
