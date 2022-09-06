-- +goose Up

INSERT INTO `sys_role`
    (`id`, `name`, `summary`)
VALUES (666, '开发组', '系统开发账号所属角色，无需单独权限授权，即可拥有系统所有权限。'),
       (10000, '超级管理员', '');

insert into `sys_module`
    (`id`, `slug`, `name`, `is_enable`, `order`)
values (1000, 'site', '站点', 1, 50);


INSERT INTO `sys_permission` (`id`, `module_id`, `parent_i1`, `parent_i2`, `name`, `slug`, `method`, `path`)
VALUES (1, 1000, 0, 0, '管理', 'manage', '', ''),
       (2, 1000, 1, 0, '权限', 'manage.permission', '', ''),
       (3, 1000, 1, 2, '创建', 'manage.permission.create', 'POST', '/admin/site/manage/permission'),
       (4, 1000, 1, 2, '修改', 'manage.permission.update', 'PUT', '/admin/site/manage/permissions/:id'),
       (5, 1000, 1, 2, '删除', 'manage.permission.delete', 'DELETE', '/admin/site/manage/permissions/:id'),
       (6, 1000, 1, 2, '列表', 'manage.permission.tree', 'GET', '/admin/site/manage/permissions'),
       (7, 1000, 1, 0, '权限', 'manage.role', '', ''),
       (8, 1000, 1, 7, '创建', 'manage.role.create', 'POST', '/admin/site/manage/role'),
       (9, 1000, 1, 7, '修改', 'manage.role.update', 'PUT', '/admin/site/manage/roles/:id'),
       (10, 1000, 1, 7, '删除', 'manage.role.delete', 'DELETE', '/admin/site/manage/roles/:id'),
       (11, 1000, 1, 7, '列表', 'manage.role.paginate', 'GET', '/admin/site/manage/roles'),
       (12, 1000, 1, 0, '账号', 'manage.admin', '', ''),
       (13, 1000, 1, 12, '创建', 'manage.admin.create', 'POST', '/admin/site/manage/admin'),
       (14, 1000, 1, 12, '修改', 'manage.admin.update', 'PUT', '/admin/site/manage/admins/:id'),
       (15, 1000, 1, 12, '删除', 'manage.admin.delete', 'DELETE', '/admin/site/manage/admins/:id'),
       (16, 1000, 1, 12, '启禁', 'manage.admin.enable', 'PUT', '/admin/site/manage/admin/enable'),
       (17, 1000, 1, 12, '列表', 'manage.admin.paginate', 'GET', '/admin/site/manage/admins'),
       (18, 1000, 0, 0, '架构', 'architecture', '', ''),
       (19, 1000, 18, 0, '模块', 'architecture.module', '', ''),
       (20, 1000, 18, 19, '创建', 'architecture.module.create', 'POST', '/admin/site/architecture/module'),
       (21, 1000, 18, 19, '修改', 'architecture.module.update', 'PUT', '/admin/site/architecture/modules/:id'),
       (22, 1000, 18, 19, '启禁', 'architecture.module.enable', 'PUT', '/admin/site/architecture/module/enable'),
       (23, 1000, 18, 19, '删除', 'architecture.module.delete', 'DELETE', '/admin/site/architecture/modules/:id'),
       (24, 1000, 18, 19, '列表', 'architecture.module.list', 'GET', '/admin/site/architecture/modules');


-- +goose Down
SELECT 'down SQL query';
