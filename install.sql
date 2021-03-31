CREATE TABLE `ccm_cluster` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `cluster_name` varchar(128) NOT NULL DEFAULT '',
  `description` varchar(512) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_all_in_one` BOOL NOT NULL DEFAULT '',
  `user` varchar(128) NOT NULL DEFAULT '',
  `status` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `cluster_name` (`cluster_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `ccm_perm` (
  `module` varchar(20) NOT NULL DEFAULT '' COMMENT '模块名',
  `action` varchar(20) NOT NULL DEFAULT '' COMMENT '操作名',
  UNIQUE KEY `module` (`module`,`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `ccm_perm` (`module`, `action`)
VALUES
	('cluster','add'),
	('cluster','del'),
	('cluster','edit'),
	('cluster','list');