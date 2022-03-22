create table dor_type_expense
(
    `id`                  int unsigned not null auto_increment,
    `type_id`             int unsigned not null comment '房型ID',
    `expense_category_id` int unsigned not null comment '费用类型ID',
    `name`                varchar(20)  not null default '' comment '名称',
    `type`                char(5)      not null default '' comment '收费：day=每日；month=每月',
    `cost`                int unsigned not null default 0 comment '费用',
    `created_at`          timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`          timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`          timestamp             default null,
    primary key (`id`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '宿舍房型费用规则表'