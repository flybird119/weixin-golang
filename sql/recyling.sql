--回收
CREATE TABLE recyling(
    --代理主键
    id uuid         primary key                                     DEFAULT gen_random_uuid(),
    store_id        TEXT                        NOT NULL,                           --店铺id
    appoint_times   JSONB                       NOT NULL            DEFAULT '[{"week":"mon","start_at":0,"end_at":24,"is_work":true},{"week":"tues","start_at":0,"end_at":24,"is_work":true},{"week":"wed","start_at":0,"end_at":24,"is_work":true},{"week":"thur","start_at":0,"end_at":24,"is_work":true},{"week":"fri","start_at":0,"end_at":24,"is_work":true},{"week":"sat","start_at":0,"end_at":24,"is_work":true},{"week":"sun","start_at":0,"end_at":24,"is_work":true}]',
    status          SMALLINT                    NOT NULL            DEFAULT 1,      --状态 1 默认收书状态  2 休息，不接受收书申请  3 账号异常，不支持一切收书行为和收书状态修改
    summary         TEXT                        NOT NULL            DEFAULT '',     --收书简介
    qrcode_url      TEXT                        NOT NULL            DEFAULT '',      --条形码url
    create_at       TIMESTAMP WITH TIME ZONE    NOT NULL            DEFAULT now(),  --更新时间
    update_at       TIMESTAMP WITH TIME ZONE    NOT NULL            DEFAULT now()   --更新时间
);
