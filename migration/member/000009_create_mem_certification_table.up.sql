create table mem_certification
(
    `id`          int unsigned not null auto_increment,
    `member_id`   varchar(64)  not null default '' comment '会员ID',
    `status`      varchar(10)  not null default '' comment '状态：auditing=待审核；valid=启用；invalid=禁用',
    `type`        varchar(20)  not null default '' comment '类型：idCard=身份证；passport=护照；other=其他',
    `name`        varchar(20)  not null default '' comment '名称',
    `no`          varchar(64)  not null default '' comment '号码',
    `other`       text         not null comment '其他',
    `valid_start` date                  default null comment '有效期：开始',
    `valid_end`   date                  default null comment '有效期：结束',
    `created_at`  timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`  timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`  timestamp             default null,
    primary key (`id`),
    key (`no`)
) auto_increment = 10000
  default collate = utf8mb4_unicode_ci comment '会员证件表'