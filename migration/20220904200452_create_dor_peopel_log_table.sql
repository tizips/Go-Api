-- +goose Up
-- +goose StatementBegin

create table dor_people_log
(
    `id`         int unsigned not null auto_increment,
    `people_id`  int unsigned not null default 0 comment '人员ID',
    `member_id`  varchar(64)  not null default '' comment '用户ID',
    `status`     varchar(5)   not null default '' comment '状态:live=入住;leave=离宿;change=调宿;refill=续住;positive=正式入住',
    `detail`     text                  default null comment '详细信息（JSON）',
    `remark`     varchar(255) not null default '' comment '备注',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default null,
    primary key (`id`),
    key (`people_id`),
    key (`member_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍人员日志表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists dor_people_log;

-- +goose StatementEnd
