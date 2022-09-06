-- +goose Up
-- +goose StatementBegin

create table sys_module
(
    `id`         int unsigned     not null auto_increment,
    `slug`       varchar(20)      not null default '' comment '标示',
    `name`       varchar(20)      not null default '' comment '名称',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `order`      int unsigned     not null default 0 comment '序号',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`)
) auto_increment = 1000
  default collate = utf8mb4_unicode_ci comment '系统模块表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists sys_module;

-- +goose StatementEnd
