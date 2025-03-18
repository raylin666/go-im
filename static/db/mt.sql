CREATE TABLE `account` (
   `id` bigint(20) NOT NULL AUTO_INCREMENT,
   `account_id` varchar(30) NOT NULL COMMENT '账号ID',
   `nickname` varchar(30) NOT NULL COMMENT '账号昵称',
   `avatar` varchar(120) NOT NULL DEFAULT '' COMMENT '账号头像',
   `is_admin` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否管理员,管理员可向任何账号发送消息 0否 1是',
   `is_online` tinyint(1) NOT NULL DEFAULT '0' COMMENT '当前状态: 0离线 1在线',
   `first_login_time` datetime DEFAULT NULL COMMENT '首次登录时间',
   `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
   `last_login_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '最后登录IP',
   `created_at` datetime NOT NULL COMMENT '创建时间',
   `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_accountid` (`account_id`) USING BTREE,
   KEY `idx_admin` (`is_admin`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账号信息表';

CREATE TABLE `account_online` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account_id` varchar(30) NOT NULL COMMENT '账号ID',
  `login_time` datetime NOT NULL COMMENT '登录时间',
  `logout_time` datetime DEFAULT NULL COMMENT '登出时间',
  `login_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '登录IP',
  `logout_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '登出IP',
  `logout_state` tinyint(1) NOT NULL DEFAULT '0' COMMENT '登出状态 0:正常退出 1:超时退出 2:服务端退出',
  `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
  `client_addr` varchar(24) NOT NULL DEFAULT '' COMMENT '客户端连接本地地址',
  `server_addr` varchar(24) NOT NULL DEFAULT '' COMMENT '服务端连接远程地址',
  `device_id` varchar(255) NOT NULL DEFAULT '' COMMENT '设备ID',
  `os` varchar(20) NOT NULL DEFAULT 'web' COMMENT '系统类型, 目前有 web|android|ios 值',
  `system` varchar(255) NOT NULL DEFAULT '' COMMENT '设备信息, 例如: HUAWEI#EML-AL00#HWEML#28#9',
  PRIMARY KEY (`id`),
  KEY `un_accountid_logintime` (`account_id`,`login_time`) USING BTREE,
  KEY `un_csaddr` (`client_addr`,`server_addr`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账号在线状态表';

CREATE TABLE `c2c_message` (
   `id` bigint(20) NOT NULL,
   `from_account` varchar(30) NOT NULL COMMENT '发送者ID',
   `to_account` varchar(30) NOT NULL COMMENT '接收者ID',
   `msg_type` tinyint(4) NOT NULL DEFAULT '3' COMMENT '消息类型 目前只支持自定义消息。1:文本 2:图片 3:自定义',
   `data` text COMMENT '消息内容',
   `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '消息状态 0:隐藏 1:显示',
   `is_revoke` tinyint(1) NOT NULL DEFAULT '0' COMMENT '消息是否已撤回 0:否 1:是',
   `revoked_at` datetime DEFAULT NULL COMMENT '消息撤回时间',
   `send_at` datetime NOT NULL COMMENT '消息发送时间',
   `from_deleted_at` datetime DEFAULT NULL COMMENT '发送者删除消息时间',
   `to_deleted_at` datetime DEFAULT NULL COMMENT '接收者删除消息时间',
   `created_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   KEY `un_from_to_account` (`from_account`,`to_account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='C2C 消息记录表';

CREATE TABLE `c2c_offline_message` (
   `id` bigint(20) NOT NULL,
   `from_account` varchar(30) NOT NULL COMMENT '发送者ID',
   `to_account` varchar(30) NOT NULL COMMENT '接收者ID',
   `message_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '消息ID 当无离线消息时为0',
   `unread_num` int(11) NOT NULL DEFAULT '0' COMMENT '未读消息数量',
   PRIMARY KEY (`id`),
   KEY `uk_from_to_account` (`from_account`,`to_account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='C2C 离线消息记录表';