CREATE TABLE IF NOT EXISTS path (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `pgroup` varchar(20) NOT NULL DEFAULT '' COMMENT 'path group',
  `path_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'path no',
  `res_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource no for the path',
  `url` varchar(128) NOT NULL DEFAULT '' COMMENT 'path url',
  `ptype` varchar(10) NOT NULL DEFAULT '' COMMENT 'path type: PROTECTED, PUBLIC',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `url` (`url`)
) ENGINE=InnoDB COMMENT='Paths';

CREATE TABLE IF NOT EXISTS resource (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `res_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource no',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource name',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `res_no` (`res_no`)
) ENGINE=InnoDB COMMENT='Resources';

CREATE TABLE IF NOT EXISTS role_resource (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  `res_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'resource no',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `role_no` (`role_no`)
) ENGINE=InnoDB COMMENT='Role resources';

CREATE TABLE IF NOT EXISTS role (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `role_no` varchar(32) NOT NULL DEFAULT '' COMMENT 'role no',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT 'name of role',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'when the record is created',
  `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who created this record',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'when the record is updated',
  `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT 'who updated this record',
  `is_del` tinyint NOT NULL DEFAULT '0' COMMENT '0-normal, 1-deleted',
  PRIMARY KEY (`id`),
  KEY `role_no` (`role_no`)
) ENGINE=InnoDB COMMENT='Roles';

-- delete from role;
-- delete from path;
-- delete from resource;
-- delete from role_resource;

-- default one for administrator
INSERT INTO role(role_no, name) VALUES ('role_554107924873216177918', 'Administrator');

INSERT INTO resource(res_no, name) VALUES
('res_556928196214785208429', 'Add Resource'),
('res_556928196214787208429', 'Add Resource To Role'),
('res_556928196214789208429', 'Remove Resource From Role'),
('res_556928196214791208429', 'Add New Role'),
('res_556928196214793208429', 'List Roles'),
('res_556928196214795208429', 'List Resources of Role'),
('res_556928196214797208429', 'List Paths'),
('res_556928196214799208429', 'Bind Path to Resource'),
('res_556928196214801208429', 'Unbind Path and Resource'),
('res_556928196214803208429', 'Delete Path'),
('res_556928196214805208429', 'Add Path'),
('res_556928196214807208429', 'Fetch Role Info');

INSERT INTO role_resource(role_no, res_no) VALUES
('role_554107924873216177918', 'res_556928196214785208429'),
('role_554107924873216177918', 'res_556928196214787208429'),
('role_554107924873216177918', 'res_556928196214789208429'),
('role_554107924873216177918', 'res_556928196214791208429'),
('role_554107924873216177918', 'res_556928196214793208429'),
('role_554107924873216177918', 'res_556928196214795208429'),
('role_554107924873216177918', 'res_556928196214797208429'),
('role_554107924873216177918', 'res_556928196214799208429'),
('role_554107924873216177918', 'res_556928196214801208429'),
('role_554107924873216177918', 'res_556928196214803208429'),
('role_554107924873216177918', 'res_556928196214805208429'),
('role_554107924873216177918', 'res_556928196214807208429');

INSERT INTO path(path_no, url, ptype, res_no, pgroup) VALUES
('path_556928196214784208429', '/open/api/resource/add', 'PROTECTED', 'res_556928196214785208429', 'goauth'),
('path_556928196214786208429', '/open/api/role/resource/add', 'PROTECTED', 'res_556928196214787208429', 'goauth'),
('path_556928196214788208429', '/open/api/role/resource/remove', 'PROTECTED', 'res_556928196214789208429', 'goauth'),
('path_556928196214790208429', '/open/api/role/add', 'PROTECTED', 'res_556928196214791208429', 'goauth'),
('path_556928196214792208429', '/open/api/role/list', 'PROTECTED', 'res_556928196214793208429', 'goauth'),
('path_556928196214794208429', '/open/api/role/resource/list', 'PROTECTED', 'res_556928196214795208429', 'goauth'),
('path_556928196214796208429', '/open/api/path/list', 'PROTECTED', 'res_556928196214797208429', 'goauth'),
('path_556928196214798208429', '/open/api/path/resource/bind', 'PROTECTED', 'res_556928196214799208429', 'goauth'),
('path_556928196214800208429', '/open/api/path/resource/unbind', 'PROTECTED', 'res_556928196214801208429', 'goauth'),
('path_556928196214802208429', '/open/api/path/delete', 'PROTECTED', 'res_556928196214803208429', 'goauth'),
('path_556928196214804208429', '/open/api/path/add', 'PROTECTED', 'res_556928196214805208429', 'goauth'),
('path_556928196214806208429', '/open/api/role/info', 'PROTECTED', 'res_556928196214807208429', 'goauth');