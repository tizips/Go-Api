create table oa_department_open
(
    `id`            int unsigned not null auto_increment,
    `department_id` int unsigned not null default 0 comment '部门ID',
    `channel`       varchar(10)  not null default '' comment '渠道：wechat=微信；dingtalk=钉钉',
    `openid`        varchar(64)  not null default '' comment '第三方平台ID',
    `created_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`    timestamp             default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment 'OA架构第三方绑定表'