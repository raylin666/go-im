CREATE TABLE `account` (
   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
   `account_id` varchar(30) NOT NULL COMMENT '账号ID',
   `nickname` varchar(30) NOT NULL DEFAULT '' COMMENT '账号昵称',
   `avatar` varchar(120) NOT NULL DEFAULT '' COMMENT '账号头像',
   `is_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否管理员,管理员可向任何账号发送消息 0否 1是',
   `status` tinyint(1) DEFAULT '0' COMMENT '当前状态: 0离线 1在线',
   `first_login_time` datetime DEFAULT NULL COMMENT '首次登录时间',
   `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
   `last_login_ip` varchar(16) NOT NULL COMMENT '最后登录IP',
   `created_at` datetime NOT NULL COMMENT '创建时间',
   `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_accountid` (`account_id`) USING BTREE,
   KEY `idx_isadmin` (`is_admin`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT "账号管理表";
