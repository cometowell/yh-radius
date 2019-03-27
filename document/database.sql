-- ----------------------------
-- Table structure for online_user
-- ----------------------------
DROP TABLE IF EXISTS `online_user`;
CREATE TABLE `online_user` (
  `id` bigint(2) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL COMMENT '账号',
  `nas_ip_addr` varchar(15) NOT NULL COMMENT 'NAS IP地址',
  `acct_session_id` varchar(128) NOT NULL COMMENT '计费session id',
  `start_time` datetime NOT NULL COMMENT '开始计费时间',
  `used_duration` int(11) NOT NULL DEFAULT '0' COMMENT '已计费时长',
  `ip_addr` varchar(15) NOT NULL COMMENT '用户IP地址',
  `mac_addr` varchar(19) NOT NULL COMMENT '用户MAC地址',
  `nas_port_id` varchar(128) DEFAULT NULL COMMENT '标识用户认证端口',
  `total_up_stream` bigint(20) NOT NULL DEFAULT '0' COMMENT '上行总流量',
  `total_down_stream` bigint(20) NOT NULL DEFAULT '0' COMMENT '下行总流量',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 COMMENT='在线用户表';

-- ----------------------------
-- Records of online_user
-- ----------------------------
INSERT INTO `online_user` VALUES ('36', 'test1', '192.168.10.236', '0001', '2019-03-08 12:20:48', '0', '192.168.10.235', '5C:FF:35:0E:58:A5', 'sdafsadfasdfasfasfdas', '0', '0');

-- ----------------------------
-- Table structure for rad_nas
-- ----------------------------
DROP TABLE IF EXISTS `rad_nas`;
CREATE TABLE `rad_nas` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `vendor_id` int(11) NOT NULL COMMENT '厂商ID',
  `name` varchar(60) NOT NULL COMMENT '名称',
  `ip_addr` varchar(15) NOT NULL COMMENT 'IP地址',
  `secret` varchar(20) NOT NULL COMMENT '共享秘钥',
  `authorize_port` int(11) NOT NULL DEFAULT '3799' COMMENT '授权端口，默认3799',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_nas_ip` (`ip_addr`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='NAS网络接入设备表';

-- ----------------------------
-- Records of rad_nas
-- ----------------------------
INSERT INTO `rad_nas` VALUES ('1', '9', '1', '192.168.10.236', '123456', '3799', '测试');
INSERT INTO `rad_nas` VALUES ('2', '0', 'test', '127.0.0.1', '123456', '3699', 'test');
INSERT INTO `rad_nas` VALUES ('99', '99', 'test', '192.168.10.233', '123456', '3799', 'mycat test');

-- ----------------------------
-- Table structure for rad_product
-- ----------------------------
DROP TABLE IF EXISTS `rad_product`;
CREATE TABLE `rad_product` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) NOT NULL COMMENT '产品名称',
  `type` int(11) NOT NULL COMMENT '类型：1:包月 2：自由时长，3：流量',
  `status` int(11) NOT NULL COMMENT '状态,0:停用，1：正常',
  `should_bind_mac_addr` tinyint(1) NOT NULL DEFAULT '0' COMMENT '需要绑定MAC地址，0：N，1：Y',
  `should_bind_vlan` tinyint(1) NOT NULL DEFAULT '0' COMMENT '需要绑定虚拟局域网，0：N，1：Y',
  `concurrent_count` int(11) NOT NULL DEFAULT '0' COMMENT '并发数',
  `product_duration` bigint(20) NOT NULL DEFAULT '0' COMMENT '时长,单位秒',
  `service_month` int(11) NOT NULL DEFAULT '0' COMMENT '套餐购买月数',
  `product_flow` bigint(20) NOT NULL DEFAULT '0' COMMENT '流量，单位KB',
  `flow_clear_cycle` tinyint(1) NOT NULL COMMENT '计费周期；0：无限时长， 1：日，2：月：3：固定（开通至使用时长截止[用户套餐过期时间]）',
  `price` int(11) NOT NULL DEFAULT '0' COMMENT '产品价格，单位分',
  `up_stream_limit` bigint(20) NOT NULL COMMENT '上行流量限制，单位Kbps',
  `down_stream_limit` bigint(20) NOT NULL COMMENT '下行流量限制,单位Kbps',
  `domain_name` varchar(200) DEFAULT NULL COMMENT '用户域',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '创建时间',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='产品表';

-- ----------------------------
-- Records of rad_product
-- ----------------------------
INSERT INTO `rad_product` VALUES ('1', '测试', '1', '1', '0', '0', '0', '0', '0', '0', '12', '0', '9000', '8000', 'aaa', '2019-03-04 17:06:51', null, null);

-- ----------------------------
-- Table structure for rad_user
-- ----------------------------
DROP TABLE IF EXISTS `rad_user`;
CREATE TABLE `rad_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` varchar(64) NOT NULL COMMENT '账号',
  `real_name` varchar(128) DEFAULT NULL COMMENT '姓名',
  `password` varchar(256) NOT NULL COMMENT '密码',
  `product_id` bigint(20) DEFAULT NULL COMMENT '产品ID',
  `status` int(11) NOT NULL COMMENT '状态，1：正常，2：停机，3：销户，4：禁用',
  `available_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '可用时长，单位：秒',
  `available_flow` bigint(20) NOT NULL DEFAULT '0' COMMENT '可用流量，单位KB',
  `expire_time` datetime DEFAULT NULL COMMENT '到期时间',
  `concurrent_count` int(11) NOT NULL DEFAULT '0' COMMENT '并发数',
  `should_bind_mac_addr` tinyint(1) NOT NULL DEFAULT '0' COMMENT '需要绑定MAC地址，0：N，1：Y',
  `should_bind_vlan` tinyint(1) NOT NULL DEFAULT '0' COMMENT '需要绑定虚拟局域网，0：N，1：Y',
  `mac_addr` varchar(19) DEFAULT NULL COMMENT 'MAC地址',
  `vlan_id` int(11) DEFAULT '0' COMMENT 'vlanId1',
  `vlan_id2` int(11) DEFAULT '0' COMMENT 'vlanId2',
  `framed_ip_addr` varchar(15) DEFAULT NULL COMMENT '用户绑定的静态IP地址',
  `installed_addr` varchar(256) DEFAULT NULL COMMENT '装机地址',
  `pause_time` datetime DEFAULT NULL COMMENT '最近停机时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_name` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Records of rad_user
-- ----------------------------
INSERT INTO `rad_user` VALUES ('1', 'test', '测试', 'oD2Ou3h126sv7bje58Z+fA==', '1', '1', '6768', '56756756', '2019-03-15 16:40:55', '1', '1', '1', '5C:FF:35:0E:58:A5', '1', '6', null, '打的费撒发送发送发大水发大厦发送', null, '2019-02-28 16:42:39', null, null);

-- ----------------------------
-- Table structure for rad_user_balance
-- ----------------------------
DROP TABLE IF EXISTS `rad_user_balance`;
CREATE TABLE `rad_user_balance` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_wallet_id` bigint(20) DEFAULT NULL COMMENT '用户钱包ID',
  `type` int(11) DEFAULT NULL COMMENT '类型1: 专项套餐，2：无限使用',
  `product_id` bigint(20) DEFAULT NULL COMMENT '产品ID',
  `balance` int(11) DEFAULT NULL COMMENT '余额',
  `expire_time` datetime DEFAULT NULL COMMENT '金额过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_wallet_id` (`user_wallet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户余额表';

-- ----------------------------
-- Records of rad_user_balance
-- ----------------------------

-- ----------------------------
-- Table structure for rad_user_special_balance
-- ----------------------------
DROP TABLE IF EXISTS `rad_user_special_balance`;
CREATE TABLE `rad_user_special_balance` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_wallet_id` bigint(20) DEFAULT NULL COMMENT '用户钱包ID',
  `type` int(11) DEFAULT NULL COMMENT '类型1: 专项套餐，2：无限使用',
  `product_id` bigint(20) DEFAULT NULL COMMENT '产品ID',
  `balance` int(11) DEFAULT NULL COMMENT '余额',
  `expire_time` datetime DEFAULT NULL COMMENT '金额过期时间',
  PRIMARY KEY (`Id`),
  UNIQUE KEY `user_wallet_id` (`user_wallet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rad_user_special_balance
-- ----------------------------

-- ----------------------------
-- Table structure for rad_user_wallet
-- ----------------------------
DROP TABLE IF EXISTS `rad_user_wallet`;
CREATE TABLE `rad_user_wallet` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户ID',
  `payment_password` varchar(256) DEFAULT NULL COMMENT '支付密码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户钱包表';

-- ----------------------------
-- Records of rad_user_wallet
-- ----------------------------

-- ----------------------------
-- Table structure for sys_department
-- ----------------------------
DROP TABLE IF EXISTS `sys_department`;
CREATE TABLE `sys_department` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `code` varchar(64) NOT NULL COMMENT '部门编码',
  `name` varchar(128) NOT NULL COMMENT '部门名称',
  `parent_id` bigint(20) NOT NULL COMMENT '上级部门ID',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='部门表';

-- ----------------------------
-- Records of sys_department
-- ----------------------------

-- ----------------------------
-- Table structure for sys_manager
-- ----------------------------
DROP TABLE IF EXISTS `sys_manager`;
CREATE TABLE `sys_manager` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `department_id` bigint(20) NOT NULL COMMENT '部门ID',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `real_name` varchar(128) DEFAULT NULL COMMENT '姓名',
  `password` varchar(256) NOT NULL COMMENT '密码',
  `status` int(11) NOT NULL COMMENT '状态，1：正常，2：停机，3：销户，4：禁用',
  `mobile` varchar(12) DEFAULT NULL COMMENT '联系方式',
  `email` varchar(250) DEFAULT NULL COMMENT '电子邮件',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_name` (`username`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统管理用户表';

-- ----------------------------
-- Records of sys_manager
-- ----------------------------

-- ----------------------------
-- Table structure for sys_manager_role_rel
-- ----------------------------
DROP TABLE IF EXISTS `sys_manager_role_rel`;
CREATE TABLE `sys_manager_role_rel` (
  `role_id` bigint(20) NOT NULL COMMENT '角色id',
  `manager_id` bigint(20) NOT NULL COMMENT '用户主键id',
  PRIMARY KEY (`role_id`,`manager_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='管理员与角色关联表';

-- ----------------------------
-- Records of sys_manager_role_rel
-- ----------------------------

-- ----------------------------
-- Table structure for sys_resource
-- ----------------------------
DROP TABLE IF EXISTS `sys_resource`;
CREATE TABLE `sys_resource` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `parent_id` bigint(20) DEFAULT NULL COMMENT '父级菜单',
  `name` varchar(255) NOT NULL COMMENT '菜单名称',
  `icon` varchar(255) DEFAULT NULL COMMENT '图标',
  `url` varchar(256) DEFAULT NULL COMMENT 'URL地址',
  `type` tinyint(1) NOT NULL COMMENT '菜单类型，0：模块，1：栏目，2：按钮',
  `enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用，1：启用，0：关闭',
  `perm_mark` varchar(255) DEFAULT NULL COMMENT '权限标志，可用于shiro注解',
  `sort_order` int(11) NOT NULL DEFAULT '1' COMMENT '排序顺序',
  `remark` varchar(512) DEFAULT NULL COMMENT '描述',
  `need_prem` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否需要权限控制,1：需要，0：不需要',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='菜单表';

-- ----------------------------
-- Records of sys_resource
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '角色名',
  `code` varchar(255) NOT NULL COMMENT '角色编码',
  `enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用角色，1：启用，0：关闭',
  `remark` varchar(512) DEFAULT NULL COMMENT '描述',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '最近更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------

-- ----------------------------
-- Table structure for sys_role_resource_rel
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_resource_rel`;
CREATE TABLE `sys_role_resource_rel` (
  `resource_id` bigint(20) NOT NULL COMMENT '菜单id',
  `role_id` bigint(20) NOT NULL COMMENT '角色id',
  PRIMARY KEY (`resource_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色与菜单关联关系表';

-- ----------------------------
-- Records of sys_role_resource_rel
-- ----------------------------

-- ----------------------------
-- Table structure for user_online_log
-- ----------------------------
DROP TABLE IF EXISTS `user_online_log`;
CREATE TABLE `user_online_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL COMMENT '账号',
  `start_time` datetime NOT NULL COMMENT '开始时间',
  `stop_time` datetime DEFAULT NULL COMMENT '结束时间',
  `used_duration` int(11) NOT NULL DEFAULT '0' COMMENT '使用时长，单位秒',
  `total_up_stream` int(11) NOT NULL DEFAULT '0' COMMENT '上行流量，单位KB',
  `total_down_stream` int(11) NOT NULL DEFAULT '0' COMMENT '下行流量，单位KB',
  `ip_addr` varchar(15) DEFAULT NULL COMMENT '用户IP地址',
  `mac_addr` varchar(17) DEFAULT NULL COMMENT '用户MAC地址',
  `nas_ip_addr` varchar(15) DEFAULT NULL COMMENT 'nas ip地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8 COMMENT='用户上网记录';

-- ----------------------------
-- Records of user_online_log
-- ----------------------------
INSERT INTO `user_online_log` VALUES ('24', 'test1', '2019-03-07 16:40:40', '2019-03-07 16:40:49', '9', '0', '0', '192.168.10.235', '5C:FF:35:0E:58:A5', '192.168.10.236');
