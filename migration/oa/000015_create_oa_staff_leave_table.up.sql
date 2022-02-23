create table oa_staff_leave
(
    `id`            int unsigned not null auto_increment,
    `member_id`     int unsigned not null default 0 comment '会员ID',
    `staff_id`      int unsigned not null default 0 comment '员工ID',
    `channel`       varchar(10)  not null default '' comment '渠道：wechat=微信；dingtalk=钉钉',
    `open_id`       int unsigned not null default 0 comment '员工开放ID',
    `last_work_day` timestamp             default null comment '最后工作日',
    `reason_type`   varchar(20)  not null default '' comment '离职原因类型：family=家庭原因；personal=个人；growing=发展；contract=合同到期不续签；relieve=协议解除；ability=无法胜任工作；layoffs=经济性裁员；regulation=严重违法违纪；other=其他',
    `reason_memo`   varchar(255) not null default '' comment '离职原因备注',
    `pre_status`    varchar(20)  not null default '' comment '离职前工作状态：try=试用期；official=正式；wait_hired=待入职',
    `status`        varchar(20)  not null default '' comment '离职状态：wait=待离职；done=已离职；no=未离职；unpass=发起离职审批但还未通过；invalid=失效（离职流程被其他流程强制终止后的状态）',
    `handover_id`   int unsigned not null default 0 comment '交接员工',
    `created_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at`    timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at`    timestamp             default null,
    primary key (`id`),
    key (`member_id`),
    key (`staff_id`),
    key (`open_id`)
) auto_increment = 1000
  default collate = utf8mb4_unicode_ci comment 'OA员工离职表'