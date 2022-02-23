INSERT INTO sys_role
    (id, name, summary)
VALUES (666, '开发组', '系统开发账号所属角色，无需单独权限授权，即可拥有系统所有权限。'),
       (10000, '超级管理员', '');

insert into sys_module
    (id, slug, name, enable, `order`)
values (1, 1, '站点', 1, 50);

INSERT INTO sys_permission
    (id, module_id, parent_i1, parent_i2, name, slug, method, path)
VALUES (1, 1, 0, 0, '授权', 'auth', '', ''),
       (2, 1, 1, 0, '权限', 'auth.permission', '', ''),
       (3, 1, 1, 2, '创建', 'auth.permission.create', 'POST', '/admin/site/management/permission'),
       (4, 1, 1, 2, '修改', 'auth.permission.update', 'PUT', '/admin/site/management/permission'),
       (5, 1, 1, 2, '删除', 'auth.permission.delete', 'DELETE', '/admin/site/management/permission'),
       (6, 1, 1, 2, '列表', 'auth.permission.tree', 'GET', '/admin/site/management/permission'),
       (7, 1, 1, 0, '权限', 'auth.role', '', ''),
       (8, 1, 1, 7, '创建', 'auth.role.create', 'POST', '/admin/site/management/role'),
       (9, 1, 1, 7, '修改', 'auth.role.update', 'PUT', '/admin/site/management/role'),
       (10, 1, 1, 7, '删除', 'auth.role.delete', 'DELETE', '/admin/site/management/role'),
       (11, 1, 1, 7, '列表', 'auth.role.paginate', 'GET', '/admin/site/management/role'),
       (12, 1, 1, 0, '账号', 'auth.admin', '', ''),
       (13, 1, 1, 12, '创建', 'auth.admin.create', 'POST', '/admin/site/management/admin'),
       (14, 1, 1, 12, '修改', 'auth.admin.update', 'PUT', '/admin/site/management/admin'),
       (15, 1, 1, 12, '删除', 'auth.admin.delete', 'DELETE', '/admin/site/management/admin'),
       (16, 1, 1, 12, '启禁', 'auth.admin.enable', 'PUT', '/admin/site/management/admin'),
       (17, 1, 1, 12, '列表', 'auth.admin.paginate', 'GET', '/admin/site/management/admin')
