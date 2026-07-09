-- SQLite 建表语句（从 MySQL 迁移）
-- 注意：SQLite 不支持 ON UPDATE CURRENT_TIMESTAMP，使用触发器替代
-- 日期列使用 DATETIME 类型，modernc.org/sqlite 驱动会自动解析为 time.Time

CREATE TABLE IF NOT EXISTS `cates` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL,
  `desc` TEXT NOT NULL,
  `domain` TEXT NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_domain` ON `cates`(`domain`);

CREATE TABLE IF NOT EXISTS `links` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL,
  `url` TEXT NOT NULL,
  `desc` TEXT NOT NULL,
  `status` TEXT NOT NULL DEFAULT 'PENDING',
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS `moods` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `content` TEXT NOT NULL DEFAULT '',
  `user_id` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS `options` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `option_key` TEXT NOT NULL,
  `option_value` TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `option_name` ON `options`(`option_key`);

CREATE TABLE IF NOT EXISTS `posts` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `cate_id` INTEGER NOT NULL DEFAULT 1,
  `type` INTEGER NOT NULL DEFAULT 1,
  `user_id` INTEGER NOT NULL,
  `title` TEXT NOT NULL DEFAULT '',
  `url` TEXT NOT NULL DEFAULT '',
  `content` TEXT NOT NULL,
  `view_num` INTEGER NOT NULL DEFAULT 1,
  `tags` TEXT,
  `status` TEXT NOT NULL DEFAULT 'ACTIVE',
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);
CREATE INDEX IF NOT EXISTS `idx_user` ON `posts`(`user_id`);
CREATE INDEX IF NOT EXISTS `idx_url` ON `posts`(`url`);
CREATE INDEX IF NOT EXISTS `idx_created_at` ON `posts`(`created_at`);
CREATE INDEX IF NOT EXISTS `idx_updated_at` ON `posts`(`updated_at`);

CREATE TABLE IF NOT EXISTS `comments` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `post_id` INTEGER NOT NULL,
  `pid` INTEGER NOT NULL,
  `name` TEXT NOT NULL DEFAULT '',
  `reply_name` TEXT NOT NULL,
  `email` TEXT NOT NULL,
  `website` TEXT NOT NULL,
  `content` TEXT NOT NULL,
  `ip` TEXT NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS `users` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL DEFAULT '',
  `password` TEXT NOT NULL DEFAULT '',
  `nick_name` TEXT NOT NULL DEFAULT '',
  `email` TEXT NOT NULL DEFAULT '',
  `status` TEXT NOT NULL DEFAULT 'ACTIVE',
  `type` INTEGER NOT NULL DEFAULT 1,
  `totp_secret` TEXT NOT NULL DEFAULT '',
  `openid` TEXT NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_user_name` ON `users`(`name`);
CREATE UNIQUE INDEX IF NOT EXISTS `un_users_email` ON `users`(`email`);

CREATE TABLE IF NOT EXISTS `reminds` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `cron` TEXT NOT NULL DEFAULT '',
  `content` TEXT NOT NULL,
  `status` TEXT NOT NULL DEFAULT 'ACTIVE',
  `next_time` DATETIME NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS `regions` (
  `region_id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `parent_id` INTEGER NOT NULL,
  `level` INTEGER NOT NULL,
  `region_name` TEXT NOT NULL DEFAULT '',
  `longitude` TEXT NOT NULL DEFAULT '',
  `latitude` TEXT NOT NULL DEFAULT '',
  `pinyin` TEXT NOT NULL,
  `az_no` TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS `idx_parent_id` ON `regions`(`parent_id`);
CREATE INDEX IF NOT EXISTS `idx_region_name` ON `regions`(`region_name`);

CREATE TABLE IF NOT EXISTS `guestbook` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL DEFAULT '',
  `content` TEXT NOT NULL,
  `ip` TEXT NOT NULL DEFAULT '',
  `top` INTEGER NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE TABLE IF NOT EXISTS `footprints` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL,
  `description` TEXT NOT NULL DEFAULT '',
  `longitude` TEXT NOT NULL,
  `latitude` TEXT NOT NULL,
  `date` TEXT NOT NULL DEFAULT '',
  `marker_color` TEXT NOT NULL DEFAULT '',
  `categories` TEXT,
  `url` TEXT NOT NULL DEFAULT '',
  `url_label` TEXT NOT NULL DEFAULT '',
  `photos` TEXT,
  `created_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
  `updated_at` DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- updated_at 自动更新触发器（替代 MySQL ON UPDATE CURRENT_TIMESTAMP）
-- PRAGMA recursive_triggers 默认 OFF，触发器内的 UPDATE 不会递归触发自身

CREATE TRIGGER IF NOT EXISTS `tr_cates_updated_at`
AFTER UPDATE ON `cates` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `cates` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;

CREATE TRIGGER IF NOT EXISTS `tr_links_updated_at`
AFTER UPDATE ON `links` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `links` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;

CREATE TRIGGER IF NOT EXISTS `tr_moods_updated_at`
AFTER UPDATE ON `moods` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `moods` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;

CREATE TRIGGER IF NOT EXISTS `tr_posts_updated_at`
AFTER UPDATE ON `posts` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `posts` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;

CREATE TRIGGER IF NOT EXISTS `tr_users_updated_at`
AFTER UPDATE ON `users` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `users` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;

CREATE TRIGGER IF NOT EXISTS `tr_reminds_updated_at`
AFTER UPDATE ON `reminds` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `reminds` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;

CREATE TRIGGER IF NOT EXISTS `tr_footprints_updated_at`
AFTER UPDATE ON `footprints` FOR EACH ROW
WHEN NEW.`updated_at` IS OLD.`updated_at`
BEGIN
    UPDATE `footprints` SET `updated_at` = datetime('now', 'localtime') WHERE `id` = NEW.`id`;
END;
