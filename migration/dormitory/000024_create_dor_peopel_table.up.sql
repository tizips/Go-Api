create table dor_people
(
    `id`          int unsigned     not null auto_increment,
    `category_id` int unsigned     not null default 0 comment '类型ID',
    `building_id` int unsigned     not null default 0 comment '楼栋ID',
    `floor_id`    int unsigned     not null default 0 comment '楼层ID',
    `room_id`     int unsigned     not null default 0 comment '房间ID',
    `bed_id`      int unsigned     not null default 0 comment '床位ID',
    `type_id`     int unsigned     not null default 0 comment '房型ID',
    `member_id`   varchar(64)      not null default '' comment '用户ID',
    `start`       timestamp                 default null comment '入住时间',
    `end`         timestamp                 default null comment '离宿时间',
    `remark`      varchar(255)     not null default '' comment '备注',
    `is_free`     tinyint unsigned not null default 0 comment '是否免费：0=否；1=是',
    `is_temp`     tinyint unsigned not null default 0 comment '是否临时：0=否；1=是',
    `is_master`   tinyint unsigned not null default 0 comment '是否责任人：0=否；1=是',
    `status`      varchar(5)       not null default '' comment '状态：live=在住；leave=离宿',
    `created_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp                 default null,
    primary key (`id`),
    key (building_id),
    key (floor_id),
    key (room_id),
    key (bed_id),
    key (type_id),
    key (member_id)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍入住表'