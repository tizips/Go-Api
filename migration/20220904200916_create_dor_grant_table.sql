-- +goose Up
-- +goose StatementBegin

create table dor_grant
(
    `id`         int unsigned not null auto_increment,
    `object`     varchar(10)  not null default '' comment '类型：package=打包；device=设备',
    `package_id` int unsigned not null default 0 comment '打包ID',
    `remark`     varchar(255) not null default '' comment '发放备注',
    `cancel`     varchar(255) not null default '' comment '取消备注',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default null,
    primary key (`id`),
    key (`package_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍发放表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_grant;

-- +goose StatementEnd
