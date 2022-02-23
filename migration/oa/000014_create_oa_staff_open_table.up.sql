create table oa_staff_open
(
    `id`         int unsigned not null auto_increment,
    `member_id`  int unsigned not null default 0 comment '会员ID',
    `staff_id`   int unsigned not null default 0 comment '员工ID',
    `channel`    varchar(10)  not null default '' comment '渠道：wechat=微信；dingtalk=钉钉',
    `openid`     varchar(64)  not null default '' comment '平台ID',
    `unionid`    varchar(64)  not null default '' comment '平台唯一标识',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default null,
    primary key (`id`)
) auto_increment = 1000
  default collate = utf8mb4_unicode_ci comment 'OA员工绑定第三方表'