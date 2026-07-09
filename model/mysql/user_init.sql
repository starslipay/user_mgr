-- Drop database if exists
DROP DATABASE IF EXISTS `user_db`;

-- Create database
CREATE DATABASE IF NOT EXISTS `user_db` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- Use database
USE `user_db`;

DROP TABLE IF EXISTS `t_relation`;
CREATE TABLE `t_relation` (
  `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
  `uid` BIGINT NOT NULL COMMENT '主键',
  `state` TINYINT NOT NULL COMMENT '关联状态',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  INDEX `idx_create_time` (`create_time`),
  INDEX `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Drop table if exists
DROP TABLE IF EXISTS `t_user_info`;
-- Create user table
CREATE TABLE `t_user_info` (
  `uid` BIGINT NOT NULL COMMENT '主键',
  `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
  `password` VARCHAR(128) NOT NULL COMMENT '支付密码',
  `name` VARCHAR(128) NOT NULL COMMENT '姓名',
  `gender` TINYINT NOT NULL COMMENT '性别',
  `age` SMALLINT NOT NULL COMMENT '年龄',
  `address` VARCHAR(128) NOT NULL COMMENT '地址',
  `phone` VARCHAR(128) NOT NULL COMMENT '手机号',
  `email` VARCHAR(128) NOT NULL COMMENT '邮箱',
  `id_type` TINYINT NOT NULL COMMENT '身份证类型',
  `id_card` VARCHAR(128) NOT NULL COMMENT '身份证号',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  INDEX `idx_create_time` (`create_time`),
  INDEX `idx_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `t_uid_segment`;
CREATE TABLE `t_uid_segment` (
  `id` BIGINT NOT NULL COMMENT '主键',
  `uid_max` BIGINT NOT NULL COMMENT '已使用的最大用户ID',
  `step` BIGINT NOT NULL COMMENT '步长',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Init test data
INSERT INTO `t_uid_segment` (`id`, `uid_max`, `step`) VALUES (1, 10000000, 1);
-- linux:  mysql -h 127.0.0.1 -P 3307 -u root -p123456 < user_init.sql
-- windows: Get-Content -Encoding UTF8 user_init.sql | mysql -h 127.0.0.1 -P 3307 -u root -p123456
-- 只读权限 multipass exec master1 -- sudo kubectl exec -it -n pay-ns mysql-0 -- mysql -ustarslipay -ppayClipayA2026
-- root权限 multipass exec master1 -- sudo kubectl exec -it -n pay-ns mysql-0 -- mysql -uroot -proot123456


-- select * from user_db.t_uid_segment limit 2;