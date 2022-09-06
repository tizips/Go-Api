-- +goose Up
-- +goose StatementBegin

create table dor_device
(
    `id`            int unsigned not null auto_increment,
    `category_id`   int unsigned not null default 0 comment '类型ID',
    `no`            varchar(64)  not null default '' comment '编号',
    `name`          varchar(64)  not null default '' comment '名称',
    `specification` varchar(64)  not null default '' comment '规格',
    `price`         int unsigned not null default 0 comment '价格',
    `unit`          varchar(64)  not null default '' comment '单位',
    `indemnity`     int unsigned not null default 0 comment '价格',
    `stock_total`   int unsigned not null default 0 comment '总库存',
    `stock_used`    int unsigned not null default 0 comment '使用量',
    `remark`        varchar(255) not null default '' comment '备注',
    `created_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`    timestamp             default null,
    primary key (`id`),
    key (`category_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍资源表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_device;

-- +goose StatementEnd
