
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

-- 为什么拆成relation和userinfo表
-- 1.职责单一，relation记录的外部和内部关联关系，userinfo记录用户的基本信息
-- 2.使用自增id作为内部主键，索引性能，分库分表都更好
-- 3.relation信息一般不会修改，可以缓存起来，通过缓存提高查询效率

-- goctl model mysql ddl -src user.sql -dir .
-- -c：开启缓存（redis，可选，不加则无缓存）
