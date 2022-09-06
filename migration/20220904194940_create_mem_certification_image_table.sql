-- +goose Up
-- +goose StatementBegin

create table mem_certification_image
(
    `id`               int unsigned not null auto_increment,
    `member_id`        varchar(64)  not null default '' comment '会员ID',
    `certification_id` int unsigned not null default 0 comment '证件ID',
    `url`              varchar(255) not null default '' comment '链接地址',
    `created_at`       timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`       timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`       timestamp             default null,
    primary key (`id`),
    key (`member_id`),
    key (`certification_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '会员证件图片表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists mem_certification_image;

-- +goose StatementEnd
