-- +goose Up
-- +goose StatementBegin

create table oa_staff_bind_department
(
    `id`            int unsigned not null auto_increment,
    `member_id`     int unsigned not null default 0 comment '会员ID',
    `staff_id`      int unsigned not null default 0 comment '员工ID',
    `department_id` int unsigned not null default 0 comment '部门ID',
    `leave_id`      int unsigned not null default 0 comment '离职ID',
    `created_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`    timestamp             default null,
    primary key (`id`),
    key (`member_id`),
    key (`staff_id`),
    key (`department_id`),
    key (`leave_id`)
) default collate = utf8mb4_unicode_ci comment 'OA员工绑定部门表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop table if exists oa_staff_bind_department;

-- +goose StatementEnd
