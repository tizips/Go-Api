create table dor_package
(
    `id`         int unsigned not null auto_increment,
    `name`       varchar(20)  not null default '' comment '名称',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍打包表'