create table dor_day
(
    `id`               int unsigned not null auto_increment,
    `category_id`      int unsigned not null default 0 comment '类型ID',
    `type_id`          int unsigned not null default 0 comment '房型ID',
    `building_id`      int unsigned not null default 0 comment '楼栋ID',
    `floor_id`         int unsigned not null default 0 comment '楼层ID',
    `room_id`          int unsigned not null default 0 comment '房间ID',
    `bed_id`           int unsigned not null default 0 comment '床位ID',
    `people_id`        int unsigned not null default 0 comment '人员ID',
    `member_id`        varchar(64)           default null comment '会员ID',
    `master_people_id` int unsigned not null default 0 comment '人员（责）ID',
    `master_member_id` varchar(64)           default null comment '会员（责）ID',
    `date`             date                  default null comment '日期',
    `created_at`       timestamp    not null default CURRENT_TIMESTAMP,
    primary key (`id`),
    key (`category_id`),
    key (`type_id`),
    key (`building_id`),
    key (`floor_id`),
    key (`room_id`),
    key (`bed_id`),
    key (`people_id`),
    key (`member_id`),
    key (`master_people_id`),
    key (`master_member_id`),
    key (`date`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍每日人员表'