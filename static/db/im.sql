CREATE TABLE `app` (
   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '应用ID',
   `ident` varchar(60) NOT NULL COMMENT '唯一标识, 用来标识来源',
   `name` varchar(30) NOT NULL COMMENT '应用名称',
   `key` bigint(20) NOT NULL COMMENT '应用KEY',
   `secret` varchar(120) NOT NULL COMMENT '应用密钥',
   `status` tinyint(1) DEFAULT '0' COMMENT '应用状态 0停用 1启用 2冻结',
   `expired_at` datetime NOT NULL COMMENT '过期时间',
   `created_at` datetime NOT NULL COMMENT '创建时间',
   `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_key` (`key`) USING BTREE,
   UNIQUE KEY `uk_ident_name` (`ident`,`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT "应用账号管理表";