--轮播图
CREATE TABLE circular(
    --代理主键
    id uuid         primary key                                     DEFAULT gen_random_uuid(),
    store_id        TEXT                        NOT NULL            DEFAULT '',     --如果为‘’那么为默认circular,每个商家初始化都要复制一遍
    type            SMALLINT                    NOT NULL            DEFAULT 1,    --轮播图类型      1 默认轮播图，不可点击  2 商品推荐  3 话题推荐
    title           TEXT                        NOT NULL            DEFAULT '', -- 来源名称 ：isbn 或者话题名称
    profile         TEXT                        NOT NULL            DEFAULT '',       --简介
    image           TEXT                        NOT NULL            DEFAULT '',         --图片url
    source_id       TEXT                        NOT NULL            DEFAULT '',
    url             TEXT                        NOT NULL            DEFAULT '',
    create_at       TIMESTAMP WITH TIME ZONE    NOT NULL            DEFAULT now(),    --更新时间
    update_at       TIMESTAMP WITH TIME ZONE    NOT NULL            DEFAULT now()    --更新时间
);
CREATE INDEX circular_index_store ON circular(store_id);
