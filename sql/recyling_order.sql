--轮播图
CREATE TABLE recyling_order(
    --代理主键
    id uuid                 primary key                                     DEFAULT gen_random_uuid(),
    store_id                TEXT                        NOT NULL,                           --店铺id
    school_id               TEXT                        NOT NULL,                           --学校id
    lp_user_id                 TEXT                     NOT NULL,                           --用户id
    images                  JSONB                                           DEFAULT '[]',   --用户上传图片
    state                   SMALLINT                    NOT NULL            DEFAULT 1,      --回收订单状态 1 默认待处理状态  2 回收搁置中  3 回收完成
    remark                  TEXT                        NOT NULL            DEFAULT '',     --用户备注
    seller_remark           TEXT                        NOT NULL            DEFAULT '',     --商家搁置备注
    addr                    TEXT                        NOT NULL          ,                 --用户地址
    mobile                  TEXT                        NOT NULL          ,                 --用户手机号
    appoint_start_at        TIMESTAMP WITH TIME ZONE    NOT NULL                      ,     --预约时间
    appoint_end_at          TIMESTAMP WITH TIME ZONE    NOT NULL                      ,     --预约时间
    create_at               TIMESTAMP WITH TIME ZONE    NOT NULL            DEFAULT now(),  --创建时间
    update_at               TIMESTAMP WITH TIME ZONE    NOT NULL            DEFAULT now()   --更新时间

);
CREATE INDEX recyling_order_index_school ON recyling_order(school_id);
