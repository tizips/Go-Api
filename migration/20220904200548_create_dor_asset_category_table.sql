-- +goose Up
-- +goose StatementBegin

create table dor_asset_category
(
    `id`         int unsigned     not null auto_increment,
    `name`       varchar(20)      not null default '' comment '名称',
    `order`      int unsigned     not null default 0 comment '序号：正序',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍资源类型表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_asset_category;

-- +goose StatementEnd
