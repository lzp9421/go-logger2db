CREATE TABLE `log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `app` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `service` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `entry` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `body` text NOT NULL,
  `create_time` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_las` (`level`,`app`,`service`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8