create table dor_grant_position
(
    `id`          int unsigned not null auto_increment,
    `grant_id`    int unsigned not null default 0 comment '发放ID',
    `object`      varchar(10)  not null default '' comment '对象：live=在住；type=房型；building=楼栋；floor=楼层；room=房间；bed=床位',
    `type_id`     int unsigned not null default 0 comment '房型ID',
    `type_bed_id` int unsigned not null default 0 comment '房型床位ID',
    `building_id` int unsigned not null default 0 comment '楼栋ID',
    `floor_id`    int unsigned not null default 0 comment '楼层ID',
    `room_id`     int unsigned not null default 0 comment '房间ID',
    `bed_id`      int unsigned not null default 0 comment '床位ID',
    `created_at`  timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp             default null,
    primary key (`id`),
    key (`grant_id`),
    key (`type_id`),
    key (`type_bed_id`),
    key (`building_id`),
    key (`floor_id`),
    key (`room_id`),
    key (`bed_id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍发放位置表'