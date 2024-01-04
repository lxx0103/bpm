-- budgets.sql
CREATE TABLE `budgets` (
    `id` int NOT NULL AUTO_INCREMENT,
    `organization_id` int NOT NULL DEFAULT 0 COMMENT '组织ID',
    `project_id` int NOT NULL DEFAULT 0 COMMENT '项目ID',
    `budget_type` tinyint NOT NULL DEFAULT 0 COMMENT '预算类型:1:材料,2:人工',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
    `quantity` int NOT NULL DEFAULT 0 COMMENT '数量',
    `unit_price` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '单价',
    `budget` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '总预算',
    `used` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '已用',
    `balance` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '剩余',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预算';

-- budget_pictures.sql
CREATE TABLE `budget_pictures` (
    `id` int NOT NULL AUTO_INCREMENT,
    `budget_id` int NOT NULL DEFAULT 0 COMMENT '预算ID',
    `link` varchar(255) NOT NULL DEFAULT '' COMMENT '图片',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='预算图片';

-- payment_requests.sql
CREATE TABLE `payment_requests` (
    `id` int NOT NULL AUTO_INCREMENT,
    `organization_id` int NOT NULL DEFAULT 0 COMMENT '组织ID',
    `project_id` int NOT NULL DEFAULT 0 COMMENT '项目ID',
    `payment_request_type` int NOT NULL DEFAULT 0 COMMENT '款项类型1采购2工款',
    `budget_id` int NOT NULL DEFAULT 0 COMMENT '预算ID',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
    `quantity` int NOT NULL DEFAULT 0 COMMENT '数量',
    `unit_price` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '单价',
    `total` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '总预算',
    `paid` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '已付',
    `due` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '剩余',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `user_id` int NOT NULL DEFAULT 0 COMMENT '用户ID',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态:1.待审核，2.审核通过，3.审核驳回，4.部分付款，5.已付款，-1.删除',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='请款记录';

-- payment_request_pictures
CREATE TABLE `payment_request_pictures` (
    `id` int NOT NULL AUTO_INCREMENT,
    `payment_request_id` int NOT NULL DEFAULT 0 COMMENT '请款记录ID',
    `link` varchar(255) NOT NULL DEFAULT '' COMMENT '图片',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='付款记录图片';

-- payment_request_history.sql
CREATE TABLE `payment_request_historys` (
    `id` int NOT NULL AUTO_INCREMENT,
    `payment_request_id` int NOT NULL DEFAULT 0 COMMENT '请款ID',
    `action` varchar(32) NOT NULL DEFAULT 0 COMMENT '操作',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='请款记录操作历史'

-- payment_records.sql
CREATE TABLE `payment_records` (
    `id` int NOT NULL AUTO_INCREMENT,
    `organization_id` int NOT NULL DEFAULT 0 COMMENT '组织ID',
    `project_id` int NOT NULL DEFAULT 0 COMMENT '项目ID',
    `payment_request_id` int NOT NULL DEFAULT 0 COMMENT '请款ID',
    `date` date NOT NULL DEFAULT '0000-00-00' COMMENT '日期',
    `amount` decimal(10,2) NOT NULL DEFAULT 0 COMMENT '金额',
    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='付款记录';

-- payment_record_pictures.sql
CREATE TABLE `payment_record_pictures` (
    `id` int NOT NULL AUTO_INCREMENT,
    `payment_record_id` int NOT NULL DEFAULT 0 COMMENT '付款记录ID',
    `link` varchar(255) NOT NULL DEFAULT '' COMMENT '图片',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='付款记录图片';

-- payment_request_type_audits.sql
CREATE TABLE `payment_request_type_audits` (
    `id` int NOT NULL AUTO_INCREMENT,
    `organization_id` int NOT NULL DEFAULT 0 COMMENT '组织ID',
    `payment_request_type` int NOT NULL DEFAULT 0 COMMENT '请款类型ID',
    `audit_level` smallint NOT NULL DEFAULT 0 COMMENT '审核级别',
    `audit_type` tinyint NOT NULL DEFAULT 0 COMMENT '审核类型：1.职位，2.用户',
    `audit_to` int NOT NULL DEFAULT 0 COMMENT '审核人/职位',
    `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态',
    `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_by` varchar(64) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='请款审核设置';

