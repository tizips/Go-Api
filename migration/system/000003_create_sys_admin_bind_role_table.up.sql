create table sys_admin_bind_role
(
    `id`         int unsigned not null auto_increment,
    `admin_id`   int unsigned not null default 0 comment '管理员ID',
    `role_id`    int unsigned not null default 0 comment '角色ID',
    `deleted_at` timestamp             default null,
    primary key (`id`),
    key (`admin_id`),
    key (`role_id`)
) default collate = utf8mb4_unicode_ci comment '系统管理员绑定角色表'