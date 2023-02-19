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
('res_578477630062593208429', 'Add Resource'),
('res_578477630062595208429', 'Add Resource To Role'),
('res_578477630062597208429', 'Remove Resource From Role'),
('res_578477630062599208429', 'Add New Role'),
('res_578477630062601208429', 'List Roles'),
('res_578477630062603208429', 'List Resources of Role'),
('res_578477630062605208429', 'List Paths'),
('res_578477630062607208429', 'Bind Path to Resource'),
('res_578477630062609208429', 'Unbind Path and Resource'),
('res_578477630062611208429', 'Delete Path'),
('res_578477630062613208429', 'Add Path'),
('res_578477630062615208429', 'Fetch Role Info'),
('res_578477630062617208429', 'Update Path Info'),
('res_585463207870465208429', 'List All Role Briefs');

INSERT INTO path(path_no, url, ptype, res_no, pgroup) VALUES
('path_578477630062592208429', '/goauth/open/api/resource/add', 'PROTECTED', 'res_578477630062593208429', 'goauth'),
('path_578477630062594208429', '/goauth/open/api/role/resource/add', 'PROTECTED', 'res_578477630062595208429', 'goauth'),
('path_578477630062596208429', '/goauth/open/api/role/resource/remove', 'PROTECTED', 'res_578477630062597208429', 'goauth'),
('path_578477630062598208429', '/goauth/open/api/role/add', 'PROTECTED', 'res_578477630062599208429', 'goauth'),
('path_578477630062600208429', '/goauth/open/api/role/list', 'PROTECTED', 'res_578477630062601208429', 'goauth'),
('path_578477630062602208429', '/goauth/open/api/role/resource/list', 'PROTECTED', 'res_578477630062603208429', 'goauth'),
('path_578477630062604208429', '/goauth/open/api/path/list', 'PROTECTED', 'res_578477630062605208429', 'goauth'),
('path_578477630062606208429', '/goauth/open/api/path/resource/bind', 'PROTECTED', 'res_578477630062607208429', 'goauth'),
('path_578477630062608208429', '/goauth/open/api/path/resource/unbind', 'PROTECTED', 'res_578477630062609208429', 'goauth'),
('path_578477630062610208429', '/goauth/open/api/path/delete', 'PROTECTED', 'res_578477630062611208429', 'goauth'),
('path_578477630062612208429', '/goauth/open/api/path/add', 'PROTECTED', 'res_578477630062613208429', 'goauth'),
('path_578477630062614208429', '/goauth/open/api/role/info', 'PROTECTED', 'res_578477630062615208429', 'goauth'),
('path_578477630062616208429', '/goauth/open/api/path/update', 'PROTECTED', 'res_578477630062617208429', 'goauth'),
('path_585463207870464208429', '/goauth/open/api/role/brief/all', 'PROTECTED', 'res_585463207870465208429', 'goauth');