CREATE TABLE `app` (
   `id` bigint(20) NOT NULL COMMENT '应用ID',
   `name` varchar(30) NOT NULL COMMENT '应用名称',
   `key` int(11) NOT NULL COMMENT '应用KEY',
   `secret` varchar(120) NOT NULL COMMENT '应用密钥',
   `status` tinyint(1) DEFAULT '0' COMMENT '应用状态 0停用 1启用 2冻结',
   `expired_at` datetime DEFAULT NULL COMMENT '过期时间',
   `created_at` datetime NOT NULL COMMENT '创建时间',
   `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_key_name` (`key`,`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT "应用账号管理表";