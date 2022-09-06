-- +goose Up
-- +goose StatementBegin

create table dor_type_bed
(
    `id`         int unsigned     not null auto_increment,
    `type_id`    int unsigned     not null default 0 comment '房型ID',
    `name`       varchar(20)      not null default '' comment '床位名称',
    `is_public`  tinyint unsigned not null default 0 comment '是否公共设备：1=是；2=否；',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`),
    key (`type_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍房型配置表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_type_bed;

-- +goose StatementEnd
