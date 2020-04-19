/* IM 数据库表 */
CREATE DATABASE thalurania;
USE thalurania;

/* 用户表 */
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
  `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id`     bigint(20) unsigned NOT NULL COMMENT '用户id',
  `nickname`    varchar(20)         NOT NULL COMMENT '昵称',
  `password`    varchar(20)         NOT NULL COMMENT '密码',
  `gender`      tinyint(4)          NOT NULL COMMENT '性别，0:未知；1:男；2:女',
  `extra`       varchar(1024)       NOT NULL COMMENT '附加属性',
  `create_time` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '用户';

/* 消息表 */
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`
(
  `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `seq`             bigint(20) unsigned NOT NULL COMMENT '消息序列号',
  `type`            tinyint(4)          NOT NULL COMMENT '消息类型',
  `content`         varchar(4094)       NOT NULL COMMENT '消息内容',
  `message_type_id` bigint(20) unsigned NOT NULL COMMENT '消息所属类型的id，1：用户；2：群组',
  `sender_type`     tinyint(3)          NOT NULL COMMENT '发送者类型，1：用户；2：系统',
  `sender_id`       bigint(20) unsigned NOT NULL COMMENT '发送者id',
  `receiver_type`   tinyint(3)          NOT NULL COMMENT '接收者类型,1:个人；2：群组',
  `receiver_id`     bigint(20) unsigned NOT NULL COMMENT '接收者id,如果是单聊信息，则为user_id，如果是群组消息，则为group_id',
  `to_user_ids`     varchar(255)        NOT NULL COMMENT '需要@的用户id列表，多个用户用，隔开',
  `send_time`       datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '消息发送时间',
  `status`          tinyint(255)        NOT NULL DEFAULT '0' COMMENT '消息状态，0：未处理；1：消息撤回',
  `create_time`     datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time`     datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_sender_id` (`sender_id`) USING BTREE,
  KEY `idx_receiver_id` (`receiver_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '消息';

/* 群组表 */
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`
(
  `id`           bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `group_id`     bigint(20)          NOT NULL COMMENT '群组id',
  `name`         varchar(50)         NOT NULL COMMENT '群组名称',
  `introduction` varchar(255) COMMENT '群组简介',
  `user_count`   int(11)             NOT NULL DEFAULT '0' COMMENT '群组人数',
  `type`         tinyint(4)          NOT NULL COMMENT '群组类型，1：小群；2：大群',
  `extra`        varchar(1024) COMMENT '附加属性',
  `create_time`  datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time`  datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT = '群组';

/* 群组成员表 */
DROP TABLE IF EXISTS `group_user`;
CREATE TABLE `group_user`
(
  `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `group_id`    bigint(20) unsigned NOT NULL COMMENT '组id',
  `user_id`     bigint(20) unsigned NOT NULL COMMENT '用户id',
  `label`       varchar(20)         NOT NULL COMMENT '用户在群组的昵称',
  `extra`       varchar(1024)       NOT NULL COMMENT '附加属性',
  `create_time` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_app_id_user_id` (`user_id`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COMMENT ='群组成员关系';