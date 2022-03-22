create table dor_room
(
    `id`          int unsigned     not null auto_increment,
    `building_id` int unsigned     not null comment '楼栋ID',
    `floor_id`    int unsigned     not null comment '楼层ID',
    `type_id`     int unsigned     not null comment '房型ID',
    `name`        varchar(20)      not null default '' comment '名称',
    `order`       int unsigned     not null default 0 comment '序号：正序',
    `is_furnish`  tinyint unsigned not null default 0 comment '是否装修：1=是；2=否；',
    `is_enable`   tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `is_public`   tinyint unsigned not null default 0 comment '是否公共区域：1=是；2=否；',
    `created_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp                 default null,
    primary key (`id`),
    key (`building_id`),
    key (`floor_id`),
    key (`type_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍房间表'