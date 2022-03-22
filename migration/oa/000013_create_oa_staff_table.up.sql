create table oa_staff
(
    `id`         int unsigned     not null auto_increment,
    `member_id`  int unsigned     not null default 0 comment '会员ID',
    `no`         varchar(255)     not null default '' comment '工号',
    `manager_id` int unsigned     not null default 0 comment '直系领导',
    `title`      varchar(255)     not null default '' comment '职位',
    `email`      varchar(255)     not null default '' comment '企业邮箱',
    `hired_date` timestamp                 default null comment '入职时间',
    `remark`     varchar(255)     not null default '' comment '备注',
    `status`     varchar(20)      not null default '' comment '状态：try=试用期；official=正式；wait_depart=待离职；depart=离职;wait_hired=待入职',
    `is_enable`  tinyint unsigned not null default 0 comment '是否启用：1=是；2=否；',
    `created_at` timestamp        not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp        not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp                 default null,
    primary key (`id`)
) auto_increment = 1000
  default collate = utf8mb4_unicode_ci comment 'OA员工表'