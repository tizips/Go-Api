create table dor_package_detail
(
    `id`         int unsigned not null auto_increment,
    `package_id` int unsigned not null default 0 comment '打包ID',
    `device_id`  int unsigned not null default 0 comment '设备ID',
    `number`     int unsigned not null default 0 comment '数量',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default null,
    primary key (`id`),
    key (`package_id`),
    key (`device_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍打包设备表'