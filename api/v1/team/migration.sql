CREATE TABLE `teams` (
  `id` int NOT NULL AUTO_INCREMENT,
  `organization_id` int NOT NULL DEFAULT '0' COMMENT '组织ID',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '班组名称',
  `leader` varchar(64) NOT NULL DEFAULT '' COMMENT '负责人',
  `phone` varchar(64) NOT NULL DEFAULT '' COMMENT '电话',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '班组状态',
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci