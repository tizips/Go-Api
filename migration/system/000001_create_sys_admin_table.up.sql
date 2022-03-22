create table sys_admin
(
    `id`         int unsigned     not null auto_increment,
    `username`   varchar(20)               default null comment '用户名',
    `mobile`     varchar(20)               default null comment '手机号',
    `email`      varchar(64)               default null comment '邮箱',
    `nickname`   varchar(32)      not null default '' comment '昵称',
    `avatar`     varchar(255)     not null default '' comment '头像',
    `password`   varchar(64)      not null default '' comment '密码',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`),
    key (`username`),
    key (`mobile`),
    key (`email`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment ='系统管理员表'