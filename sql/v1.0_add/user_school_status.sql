-- 中国高校所有专业名称
CREATE TABLE IF NOT EXISTS user_school_status (

    id                      UUID PRIMARY KEY         NOT NULL                 DEFAULT gen_random_uuid(),
    school_id               TEXT                      NOT NULL                          ,
    user_id                 TEXT                      NOT NULL                          ,
    institute_id            TEXT                      NOT NULL                          ,
    institute_major_id     TEXT                      NOT NULL                           ,
    create_at               TIMESTAMP WITH TIME ZONE  NOT NULL                 DEFAULT now(),   --创建时间
    update_at               TIMESTAMP WITH TIME ZONE  NOT NULL                 DEFAULT now()    --更新时间
);
