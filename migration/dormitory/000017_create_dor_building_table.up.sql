create table dor_building
(
    `id`         int unsigned     not null auto_increment,
    `name`       varchar(20)      not null default '' comment '名称',
    `order`      int unsigned     not null default 0 comment '序号：正序',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `is_public`  tinyint unsigned not null default 0 comment '是否公共区域：1=是；2=否；',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍楼栋表'