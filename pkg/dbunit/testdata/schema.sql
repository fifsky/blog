-- SQLite 建表语句（dbunit 自身测试用）

CREATE TABLE IF NOT EXISTS `actions` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `user_id` INTEGER NOT NULL,
  `content` TEXT NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS `articles` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `pid` INTEGER NOT NULL,
  `doc_id` INTEGER NOT NULL,
  `title` TEXT NOT NULL DEFAULT '',
  `user_id` INTEGER NOT NULL,
  `last_user_id` INTEGER NOT NULL,
  `content` TEXT NOT NULL,
  `type` INTEGER NOT NULL,
  `link` TEXT NOT NULL DEFAULT '',
  `reference_id` INTEGER NOT NULL DEFAULT 0,
  `edit_type` INTEGER NOT NULL DEFAULT 1,
  `status` INTEGER NOT NULL DEFAULT 1,
  `sort` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS `idx_created_at` ON `articles`(`created_at`);
CREATE INDEX IF NOT EXISTS `idx_updated_at` ON `articles`(`updated_at`);

CREATE TABLE IF NOT EXISTS `documents` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `user_id` INTEGER NOT NULL,
  `last_user_id` INTEGER NOT NULL,
  `title` TEXT NOT NULL DEFAULT '',
  `domain` TEXT NOT NULL DEFAULT '',
  `logo` TEXT NOT NULL DEFAULT '',
  `description` TEXT NOT NULL DEFAULT '',
  `permission` INTEGER NOT NULL DEFAULT 1,
  `password` TEXT NOT NULL DEFAULT '',
  `status` INTEGER NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_domain` ON `documents`(`domain`);
CREATE INDEX IF NOT EXISTS `idx_doc_created_at` ON `documents`(`created_at`);
CREATE INDEX IF NOT EXISTS `idx_doc_updated_at` ON `documents`(`updated_at`);

CREATE TABLE IF NOT EXISTS `histories` (
  `hid` INTEGER PRIMARY KEY AUTOINCREMENT,
  `id` INTEGER NOT NULL,
  `pid` INTEGER NOT NULL,
  `doc_id` INTEGER NOT NULL,
  `title` TEXT NOT NULL DEFAULT '',
  `user_id` INTEGER NOT NULL,
  `last_user_id` INTEGER NOT NULL,
  `content` TEXT NOT NULL,
  `type` INTEGER NOT NULL,
  `link` TEXT NOT NULL DEFAULT '',
  `reference_id` INTEGER NOT NULL,
  `edit_type` INTEGER NOT NULL,
  `status` INTEGER NOT NULL,
  `sort` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS `members` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `doc_id` INTEGER NOT NULL,
  `user_id` INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_doc_user` ON `members`(`doc_id`, `user_id`);

CREATE TABLE IF NOT EXISTS `shares` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL DEFAULT '',
  `domain` TEXT NOT NULL DEFAULT '',
  `doc_id` INTEGER NOT NULL,
  `password` TEXT NOT NULL DEFAULT '',
  `share_ids` TEXT NOT NULL DEFAULT '',
  `user_id` INTEGER NOT NULL,
  `status` INTEGER NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_share_domain` ON `shares`(`domain`);

CREATE TABLE IF NOT EXISTS `users` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `user_name` TEXT NOT NULL DEFAULT '',
  `email` TEXT NOT NULL DEFAULT '',
  `real_name` TEXT NOT NULL DEFAULT '',
  `password` TEXT NOT NULL DEFAULT '',
  `avatar` TEXT NOT NULL DEFAULT '',
  `status` INTEGER NOT NULL DEFAULT 1,
  `about` TEXT NOT NULL DEFAULT '',
  `role` TEXT NOT NULL DEFAULT 'user',
  `organization` TEXT NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_email` ON `users`(`email`);
CREATE UNIQUE INDEX IF NOT EXISTS `un_user_name` ON `users`(`user_name`);

CREATE TABLE IF NOT EXISTS `custom` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL,
  `nick_name` TEXT NOT NULL,
  `status` INTEGER NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS `un_name` ON `custom`(`name`);
CREATE UNIQUE INDEX IF NOT EXISTS `un_nick_name` ON `custom`(`nick_name`);
