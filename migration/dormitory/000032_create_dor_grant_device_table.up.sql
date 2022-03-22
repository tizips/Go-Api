create table dor_grant_device
(
    `id`         int unsigned not null auto_increment,
    `grant_id`   int unsigned not null default 0 comment '发放ID',
    `device_id`  int unsigned not null default 0 comment '设备ID',
    `number`     int unsigned not null default 0 comment '数量',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default null,
    primary key (`id`),
    key (`grant_id`),
    key (`device_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍发放设备表'