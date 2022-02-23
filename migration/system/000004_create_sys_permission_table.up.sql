CREATE TABLE sys_permission
(
    `id`         int unsigned not null auto_increment,
    `module_id`  int(20)      not null default 0 comment '模块ID',
    `parent_i1`  int unsigned not null default 0 comment '父级1',
    `parent_i2`  int unsigned not null default 0 comment '父级2',
    `name`       varchar(20)  not null default '' comment '名称',
    `slug`       varchar(64)  not null default '' comment '标示',
    `method`     varchar(10)  not null default '' comment 'method',
    `path`       varchar(64)  not null default '' comment 'path',
    `created_at` timestamp    not null default CURRENT_TIMESTAMP,
    `updated_at` timestamp    not null default CURRENT_TIMESTAMP,
    `deleted_at` timestamp             default NULL,
    primary key (`id`),
    key (`module_id`),
    key (`parent_i1`),
    key (`parent_i2`)
) auto_increment = 10000
  default collate utf8mb4_unicode_ci comment '系统权限表'