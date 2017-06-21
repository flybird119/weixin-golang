create sequence goods_batch_upload_id_seq;

create table goods_batch_upload(

    id text primary key not null default trim(to_char(nextval('goods_batch_upload_id_seq'),'000000000')),--id
    store_id        text            not null,                                 --云店铺id
    success_num     int             not null            default 0,                 --成功数量
    failed_num      int             not null            default 0,                 --失败数量
    state           int             not null            default 1,                 --状态 1 导入中 2 失败  3 导入完成
    type            int             not null            default 0,                 --上传类型 0新书 1 旧书
    discount        int             not null            ,                          --折扣
    storehouse_id   text            not null            default '',                --仓库id
    shelf_id        text            not null            default '',                --货架id
    floor_id        text            not null            default '',                --货架层id
    error_reason    text            not null            default '',                --失败原因 配合 状态2 使用
    origin_file     text            not null            default '',                --商家源文件
    origin_filename text            not null,
    error_file      text            not null            default '',                --处理错误文件
    create_at       timestamptz     not null            default now(),             --创建时间
    complete_at     timestamptz     ,                                              --导入完成时间
    update_at       timestamptz     not null            default now()              --更改时间
);
