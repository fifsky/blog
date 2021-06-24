# Dump of table cates
# ------------------------------------------------------------

DROP TABLE IF EXISTS `cates`;

CREATE TABLE `cates` (
                         `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                         `name` varchar(100) NOT NULL DEFAULT '',
                         `desc` varchar(255) NOT NULL DEFAULT '',
                         `domain` varchar(100) NOT NULL DEFAULT '',
                         `created_at` datetime NOT NULL,
                         `updated_at` datetime NOT NULL,
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `un_domain` (`domain`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

# Dump of table links
# ------------------------------------------------------------

DROP TABLE IF EXISTS `links`;

CREATE TABLE `links` (
                         `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                         `name` varchar(100) NOT NULL DEFAULT '',
                         `url` varchar(200) NOT NULL DEFAULT '',
                         `desc` varchar(255) NOT NULL DEFAULT '',
                         `created_at` datetime NOT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

# Dump of table moods
# ------------------------------------------------------------

DROP TABLE IF EXISTS `moods`;

CREATE TABLE `moods` (
                         `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                         `content` varchar(255) NOT NULL DEFAULT '',
                         `user_id` int(10) unsigned NOT NULL,
                         `created_at` datetime NOT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table options
# ------------------------------------------------------------

DROP TABLE IF EXISTS `options`;

CREATE TABLE `options` (
                           `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
                           `option_key` varchar(100) NOT NULL DEFAULT '',
                           `option_value` varchar(200) NOT NULL DEFAULT '',
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `option_name` (`option_key`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table posts
# ------------------------------------------------------------

DROP TABLE IF EXISTS `posts`;

CREATE TABLE `posts` (
                         `id` int unsigned NOT NULL AUTO_INCREMENT,
                         `cate_id` int unsigned NOT NULL DEFAULT '1',
                         `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:文章,2:页面',
                         `user_id` int unsigned NOT NULL,
                         `title` varchar(200) NOT NULL DEFAULT '',
                         `url` varchar(100) NOT NULL DEFAULT '' COMMENT '页面缩略名',
                         `content` longtext NOT NULL,
                         `view_num` int NOT NULL DEFAULT '1',
                         `created_at` datetime NOT NULL COMMENT '创建时间',
                         `updated_at` datetime NOT NULL COMMENT '更新时间',
                         `status` int NOT NULL DEFAULT '1' COMMENT '状态',
                         PRIMARY KEY (`id`),
                         KEY `post_author` (`user_id`),
                         KEY `type_status_date` (`id`,`type`),
                         KEY `post_name` (`url`) USING BTREE,
                         KEY `post_title` (`title`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table comments
# ------------------------------------------------------------
DROP TABLE IF EXISTS `comments`;

CREATE TABLE `comments` (
                            `id` int unsigned NOT NULL AUTO_INCREMENT,
                            `post_id` int NOT NULL COMMENT '文章PID',
                            `pid` int NOT NULL COMMENT '回复评论ID',
                            `name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
                            `content` tinytext NOT NULL COMMENT '内容',
                            `ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'IP',
                            `created_at` datetime NOT NULL COMMENT '评论时间',
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
                         `id` int unsigned NOT NULL AUTO_INCREMENT,
                         `name` varchar(100) NOT NULL DEFAULT '',
                         `password` varchar(100) NOT NULL DEFAULT '',
                         `nick_name` varchar(100) NOT NULL DEFAULT '',
                         `email` varchar(100) NOT NULL DEFAULT '',
                         `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1正常，2删除',
                         `type` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '1:管理员,2:编辑',
                         `created_at` datetime NOT NULL,
                         `updated_at` datetime NOT NULL,
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `un_user_name` (`name`),
                         UNIQUE KEY `uix_users_email` (`email`),
                         UNIQUE KEY `uix_users_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table reminds
# ------------------------------------------------------------

DROP TABLE IF EXISTS `reminds`;

CREATE TABLE `reminds` (
                           `id` int unsigned NOT NULL AUTO_INCREMENT,
                           `type` int NOT NULL COMMENT '0固定，1每分钟，2每个小时，3每周，4，每天，5，每月，6，每年',
                           `content` varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
                           `month` int NOT NULL COMMENT '月',
                           `week` int NOT NULL COMMENT '周',
                           `day` int NOT NULL COMMENT '日',
                           `hour` int NOT NULL COMMENT '时',
                           `minute` int NOT NULL COMMENT '分',
                           `created_at` datetime NOT NULL COMMENT '创建时间',
                           `status` int NOT NULL DEFAULT '1',
                           `next_time` datetime NOT NULL COMMENT '下次提醒时间',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;