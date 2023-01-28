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

-- default one for administrator
INSERT INTO role(role_no, name) VALUES ('role_554107924873216177918', 'Administrator');

INSERT INTO resource(res_no, name) VALUES
  ('res_555442491572225208429', 'Add Resource'),
  ('res_555442491572227208429', 'Add Resource To Role'),
  ('res_555442491572229208429', 'Remove Resource From Role'),
  ('res_555442491572231208429', 'Add New Role'),
  ('res_555442491572233208429', 'List Roles'),
  ('res_555442491572235208429', 'List Resources of Role'),
  ('res_555442491572237208429', 'List Paths'),
  ('res_555442491572239208429', 'Bind Path to Resource'),
  ('res_555442491572241208429', 'Unbind Path and Resource'),
  ('res_555442491572243208429', 'Delete Path'),
  ('res_555442491572245208429', 'Add Path');

INSERT INTO role_resource(role_no, res_no) VALUES
  ('role_554107924873216177918', 'res_555442491572225208429'),
  ('role_554107924873216177918', 'res_555442491572227208429'),
  ('role_554107924873216177918', 'res_555442491572229208429'),
  ('role_554107924873216177918', 'res_555442491572231208429'),
  ('role_554107924873216177918', 'res_555442491572233208429'),
  ('role_554107924873216177918', 'res_555442491572235208429'),
  ('role_554107924873216177918', 'res_555442491572237208429'),
  ('role_554107924873216177918', 'res_555442491572239208429'),
  ('role_554107924873216177918', 'res_555442491572241208429'),
  ('role_554107924873216177918', 'res_555442491572243208429'),
  ('role_554107924873216177918', 'res_555442491572245208429');

INSERT INTO path(path_no, url, ptype, res_no, pgroup) VALUES
  ('path_555442491572224208429', '/open/api/resource/add', 'PROTECTED', 'res_555442491572225208429', 'goauth'),
  ('path_555442491572226208429', '/open/api/role/resource/add', 'PROTECTED', 'res_555442491572227208429', 'goauth'),
  ('path_555442491572228208429', '/open/api/role/resource/remove', 'PROTECTED', 'res_555442491572229208429', 'goauth'),
  ('path_555442491572230208429', '/open/api/role/add', 'PROTECTED', 'res_555442491572231208429', 'goauth'),
  ('path_555442491572232208429', '/open/api/role/list', 'PROTECTED', 'res_555442491572233208429', 'goauth'),
  ('path_555442491572234208429', '/open/api/role/resource/list', 'PROTECTED', 'res_555442491572235208429', 'goauth'),
  ('path_555442491572236208429', '/open/api/path/list', 'PROTECTED', 'res_555442491572237208429', 'goauth'),
  ('path_555442491572238208429', '/open/api/path/resource/bind', 'PROTECTED', 'res_555442491572239208429', 'goauth'),
  ('path_555442491572240208429', '/open/api/path/resource/unbind', 'PROTECTED', 'res_555442491572241208429', 'goauth'),
  ('path_555442491572242208429', '/open/api/path/delete', 'PROTECTED', 'res_555442491572243208429', 'goauth'),
  ('path_555442491572244208429', '/open/api/path/add', 'PROTECTED', 'res_555442491572245208429', 'goauth');