/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50721
Source Host           : localhost:3306
Source Database       : radius-v2

Target Server Type    : MYSQL
Target Server Version : 50721
File Encoding         : 65001

Date: 2019-07-08 15:20:45
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for rad_area
-- ----------------------------
DROP TABLE IF EXISTS `rad_area`;
CREATE TABLE `rad_area` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '编码',
  `name` varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '大区名',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态,1：正常，2：停用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `description` varchar(1000) COLLATE utf8mb4_bin DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户片区划分表';

-- ----------------------------
-- Records of rad_area
-- ----------------------------
INSERT INTO `rad_area` VALUES ('1', 'test', '测试片区', '1', '2019-07-04 15:44:06', '2019-07-05 10:52:56', '测试片区');

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
) ENGINE=InnoDB AUTO_INCREMENT=102 DEFAULT CHARSET=utf8 COMMENT='NAS网络接入设备表';

-- ----------------------------
-- Records of rad_nas
-- ----------------------------
INSERT INTO `rad_nas` VALUES ('2', '9', 'test', '127.0.0.1', '123456', '3699', 'test');
INSERT INTO `rad_nas` VALUES ('99', '2011', 'test2', '10.18.10.68', '123456', '3799', 'mycat test');

-- ----------------------------
-- Table structure for rad_online_user
-- ----------------------------
DROP TABLE IF EXISTS `rad_online_user`;
CREATE TABLE `rad_online_user` (
  `id` bigint(2) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL COMMENT '账号',
  `nas_ip_addr` varchar(15) NOT NULL COMMENT 'NAS IP地址',
  `acct_session_id` varchar(128) NOT NULL COMMENT '计费session id',
  `start_time` datetime NOT NULL COMMENT '开始计费时间',
  `used_duration` int(11) NOT NULL DEFAULT '0' COMMENT '已计费时长',
  `ip_addr` varchar(15) NOT NULL COMMENT '用户IP地址',
  `mac_addr` varchar(19) NOT NULL COMMENT '用户MAC地址',
  `nas_port_id` varchar(128) DEFAULT NULL COMMENT '标识用户认证端口',
  `total_up_stream` bigint(20) NOT NULL DEFAULT '0' COMMENT '上行总流量 KB',
  `total_down_stream` bigint(20) NOT NULL DEFAULT '0' COMMENT '下行总流量，KB',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8 COMMENT='在线用户表';

-- ----------------------------
-- Records of rad_online_user
-- ----------------------------

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
  `flow_clear_cycle` tinyint(1) NOT NULL COMMENT '计费周期；1：默认， 2：日，3：月：4：固定（开通至使用时长截止[用户套餐过期时间]）',
  `price` int(11) NOT NULL DEFAULT '0' COMMENT '产品价格，单位分',
  `up_stream_limit` bigint(20) NOT NULL COMMENT '上行流量限制，单位Mbps',
  `down_stream_limit` bigint(20) NOT NULL COMMENT '下行流量限制,单位Mbps',
  `domain_name` varchar(200) DEFAULT NULL COMMENT '用户域',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '创建时间',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COMMENT='产品表';

-- ----------------------------
-- Records of rad_product
-- ----------------------------

-- ----------------------------
-- Table structure for rad_town
-- ----------------------------
DROP TABLE IF EXISTS `rad_town`;
CREATE TABLE `rad_town` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `area_id` bigint(20) NOT NULL COMMENT '片区ID',
  `code` varchar(100) COLLATE utf8mb4_bin NOT NULL COMMENT '编码',
  `name` varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '大区名',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态,1：正常，2：停用',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(1000) COLLATE utf8mb4_bin DEFAULT '' COMMENT '描述',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='村镇/街道，片区下级单位表';

-- ----------------------------
-- Records of rad_town
-- ----------------------------
INSERT INTO `rad_town` VALUES ('1', '1', 'test', '测试街道', '1', '2019-07-04 15:43:21', '2019-07-05 15:02:57', '测试街道');

-- ----------------------------
-- Table structure for rad_user
-- ----------------------------
DROP TABLE IF EXISTS `rad_user`;
CREATE TABLE `rad_user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` varchar(64) NOT NULL COMMENT '账号',
  `real_name` varchar(128) DEFAULT NULL COMMENT '姓名',
  `town_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '村镇/街道ID',
  `password` varchar(256) NOT NULL COMMENT '密码',
  `product_id` bigint(20) DEFAULT NULL COMMENT '产品ID',
  `status` int(11) NOT NULL COMMENT '状态，1：正常，2：停机，3：禁用，4：销户',
  `available_time` bigint(20) NOT NULL DEFAULT '0' COMMENT '可用时长，单位：秒',
  `available_flow` bigint(20) NOT NULL DEFAULT '0' COMMENT '可用流量，单位KB',
  `expire_time` date DEFAULT NULL COMMENT '到期时间',
  `concurrent_count` int(11) NOT NULL DEFAULT '0' COMMENT '并发数',
  `should_bind_mac_addr` tinyint(1) NOT NULL DEFAULT '0' COMMENT '需要绑定MAC地址，1：Y，2：N',
  `should_bind_vlan` tinyint(1) NOT NULL DEFAULT '0' COMMENT '需要绑定虚拟局域网，1：Y，2：N',
  `mac_addr` varchar(19) DEFAULT NULL COMMENT 'MAC地址',
  `vlan_id` int(11) DEFAULT '0' COMMENT 'vlanId1',
  `vlan_id2` int(11) DEFAULT '0' COMMENT 'vlanId2',
  `framed_ip_addr` varchar(15) DEFAULT NULL COMMENT '用户绑定的静态IP地址',
  `installed_addr` varchar(256) DEFAULT NULL COMMENT '装机地址',
  `mobile` varchar(12) DEFAULT NULL COMMENT '手机号码',
  `email` varchar(200) DEFAULT NULL COMMENT '电子邮件',
  `pause_time` datetime DEFAULT NULL COMMENT '最近停机时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_name` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Records of rad_user
-- ----------------------------

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
-- Table structure for rad_user_online_log
-- ----------------------------
DROP TABLE IF EXISTS `rad_user_online_log`;
CREATE TABLE `rad_user_online_log` (
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
-- Records of rad_user_online_log
-- ----------------------------

-- ----------------------------
-- Table structure for rad_user_order_record
-- ----------------------------
DROP TABLE IF EXISTS `rad_user_order_record`;
CREATE TABLE `rad_user_order_record` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '用户id',
  `product_id` bigint(20) NOT NULL COMMENT '产品id',
  `price` int(11) NOT NULL COMMENT '价格，单位：分',
  `sys_user_id` bigint(20) NOT NULL COMMENT '操作管理员',
  `order_time` datetime NOT NULL COMMENT '订单时间',
  `status` tinyint(1) NOT NULL DEFAULT '2' COMMENT '1:预定，2: 已生效，3：已取消',
  `end_date` date NOT NULL COMMENT '订单截止日期',
  `count` int(11) NOT NULL DEFAULT '1' COMMENT '套餐倍数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8 COMMENT='用户订单表';

-- ----------------------------
-- Records of rad_user_order_record
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
  `status` int(1) NOT NULL DEFAULT '1' COMMENT '1：正常，2：停用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COMMENT='部门表';

-- ----------------------------
-- Records of sys_department
-- ----------------------------
INSERT INTO `sys_department` VALUES ('1', 'test', '测试', '0', '2019-07-03 15:53:29', null, '发发发22', '1');

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
  `type` tinyint(1) NOT NULL COMMENT '菜单类型，1：模块，2：栏目，3：按钮',
  `enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用，1：启用，0：关闭',
  `perm_mark` varchar(255) DEFAULT NULL COMMENT '权限标志，可用于shiro注解',
  `sort_order` int(11) NOT NULL DEFAULT '1' COMMENT '排序顺序',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  `should_perm_control` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否需要权限控制,1：需要，0：不需要',
  `level` tinyint(1) NOT NULL COMMENT '层次',
  `front_router` varchar(200) DEFAULT NULL COMMENT '前端路由',
  `front_key` varchar(255) DEFAULT NULL COMMENT '前端路由key',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=476 DEFAULT CHARSET=utf8 COMMENT='菜单表';

-- ----------------------------
-- Records of sys_resource
-- ----------------------------
INSERT INTO `sys_resource` VALUES ('100', '0', '用户管理', 'team', '/user/list', '2', '1', 'user::list', '100', '用户管理', '1', '1', '/user', 'user');
INSERT INTO `sys_resource` VALUES ('110', '100', '添加用户', '', '/user/add', '3', '1', 'user::add', '110', '添加用户', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('120', '100', '修改用户', '', '/user/update', '3', '1', 'user::list', '120', '修改用户', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('130', '100', '删除用户', null, '/user/delete', '3', '1', 'user::delete', '130', '删除用户', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('140', '100', '获取用户信息', '', '/user/info', '3', '1', 'user::info', '130', '获取用户信息', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('150', '100', '用户续订', null, '/user/continue', '3', '1', 'user::continue', '1', '用户续订', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('160', '100', '获取用户订购记录', null, '/user/order/record', '3', '1', 'user::order::record', '160', '获取用户订购记录', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('200', '0', '套餐管理', 'shopping', '/product/list', '2', '1', 'product::list', '200', '套餐管理', '1', '1', '/product', 'product');
INSERT INTO `sys_resource` VALUES ('210', '200', '添加套餐', '', '/product/add', '3', '1', 'product::add', '200', '添加套餐', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('220', '200', '套餐信息', '', '/product/info', '3', '1', 'product::add', '200', '添加套餐', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('230', '200', '修改套餐', '', '/product/update', '3', '1', 'product::add', '200', '添加套餐', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('240', '200', '删除套餐', '', '/product/delete', '3', '1', 'product::add', '200', '添加套餐', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('250', '200', '获取套餐列表', '', '/fetch/product', '3', '1', 'product::fetch', '200', '获取套餐列表', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('300', '0', '在线用户', 'global', '/online/list', '2', '1', 'online::list', '300', '在线用户', '1', '1', '/online', 'online');
INSERT INTO `sys_resource` VALUES ('310', '300', '用户下线', '', '/online/off', '2', '1', 'online::off', '310', '用户下线', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('320', '300', '清理在线用户', null, '/online/delete', '3', '1', 'online::delete', '300', '清理在线用户', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('400', '0', '系统设置', 'setting', '', '1', '1', '', '400', '系统设置', '1', '1', null, 'system');
INSERT INTO `sys_resource` VALUES ('410', '400', '管理员', 'user', '/system/user/list', '2', '1', 'manager::list', '410', '管理员', '1', '2', '/sysUser', 'manager');
INSERT INTO `sys_resource` VALUES ('411', '410', '添加管理员', null, '/system/user/add', '3', '1', 'manager::add', '410', '添加管理员', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('412', '410', '修改管理员', null, '/system/user/update', '3', '1', 'manager::update', '410', '修改管理员', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('413', '410', '删除管理员', null, '/system/user/delete', '3', '1', 'manager::delete', '410', '删除管理员', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('414', '410', '获取管理员信息', null, '/system/user/info', '3', '1', 'manager::info', '410', '获取管理员信息', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('415', '410', '修改管理员密码', null, '/system/user/change/password', '3', '1', '/manager::change::password', '410', '修改管理员密码', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('416', '410', '获取会话用户信息', null, '/system/user/session/info', '3', '1', 'system::user::session::info', '416', null, '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('420', '400', 'NAS管理', 'database', '/nas/list', '2', '1', 'nas::list', '420', 'NAS管理', '1', '2', '/nas', 'nas');
INSERT INTO `sys_resource` VALUES ('421', '420', '添加NAS', null, '/nas/add', '3', '1', 'nas::add', '420', '添加NAS', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('422', '420', '修改NAS', null, '/nas/update', '3', '1', 'nas::update', '420', '修改NAS', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('423', '420', '删除NAS', null, '/nas/delete', '3', '1', 'nas::delete', '420', '删除NAS', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('424', '420', '获取NAS信息', null, '/nas/info', '3', '1', 'nas::info', '420', '获取NAS信息', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('430', '400', '部门管理', 'appstore', '/department/list', '2', '1', 'department::list', '430', '部门管理', '1', '2', '/department', 'department');
INSERT INTO `sys_resource` VALUES ('431', '430', '添加部门', null, '/department/add', '3', '1', 'department::add', '430', '添加部门', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('432', '430', '修改部门', null, '/department/update', '3', '1', 'department::update', '430', '修改部门', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('433', '430', '删除部门', null, '/department/delete', '3', '1', 'department::delete', '430', '删除部门', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('434', '430', '获取部门列表', null, '/fetch/department', '3', '1', 'department::fetch', '430', '获取部门列表', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('435', '430', '部门信息', null, '/department/info', '2', '1', 'department::info', '430', '部门信息', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('440', '400', '角色管理', 'solution', '/role/list', '2', '1', 'role::list', '440', '角色管理', '1', '2', '/role', 'role');
INSERT INTO `sys_resource` VALUES ('441', '440', '添加角色', null, '/role/add', '3', '1', 'role::add', '440', '添加角色', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('442', '440', '获取角色信息', null, '/role/info', '3', '1', 'role::info', '440', '获取角色信息', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('443', '440', '修改角色', null, '/role/update', '3', '1', 'role::update', '440', '修改角色', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('444', '440', '删除角色', null, '/role/delete', '3', '1', 'role::delete', '440', '删除角色', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('445', '440', '角色赋权', null, '/role/empower/\\d+', '3', '1', 'role::empower', '440', '角色赋权', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('446', '440', '获取角色权限', null, '/role/resources', '3', '1', 'role::resource', '440', '获取角色权限', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('447', '440', '构建菜单', null, '/session/resource', '3', '1', 'session::resource', '440', '构建菜单', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('450', '400', '菜单管理', 'profile', '/resource/list', '2', '1', 'resource::list', '450', '菜单管理', '1', '2', '/resource', 'resource');
INSERT INTO `sys_resource` VALUES ('460', '400', '片区管理', null, '/area/list', '2', '1', 'area::list', '460', '部门管理', '1', '2', '/area', 'area');
INSERT INTO `sys_resource` VALUES ('461', '460', '添加片区', null, '/area/add', '3', '1', 'area::add', '461', '添加片区', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('462', '460', '修改片区', null, '/area/update', '3', '1', 'area::update', '462', '修改片区', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('463', '460', '删除片区', null, '/area/delete', '3', '1', 'area::delete', '463', '删除片区', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('464', '460', '获取片区列表', null, '/fetch/areas', '3', '1', 'area::fetch', '464', '获取片区列表', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('465', '460', '片区信息', null, '/area/info', '2', '1', 'area::info', '465', '片区信息', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('470', '400', '村镇街道管理', null, '/town/list', '2', '1', 'town::list', '470', '村镇街道管理', '1', '2', '/town', 'town');
INSERT INTO `sys_resource` VALUES ('471', '470', '添加村镇街道', null, '/town/add', '3', '1', 'town::add', '471', '添加村镇街道', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('472', '470', '修改村镇街道', null, '/town/update', '3', '1', 'town::update', '472', '修改村镇街道', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('473', '470', '删除村镇街道', null, '/town/delete', '3', '1', 'town::delete', '473', '删除村镇街道', '1', '3', null, null);
INSERT INTO `sys_resource` VALUES ('474', '470', '获取村镇街道列表', null, '/fetch/towns', '3', '1', 'town::fetch', '474', '获取村镇街道列表', '0', '3', null, null);
INSERT INTO `sys_resource` VALUES ('475', '470', '村镇街道信息', null, '/town/info', '2', '1', 'town::info', '475', '村镇街道信息', '0', '3', null, null);

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '角色名',
  `code` varchar(255) NOT NULL COMMENT '角色编码',
  `enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用角色，1：启用，2：关闭',
  `description` varchar(512) DEFAULT NULL COMMENT '描述',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '最近更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='角色表';

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES ('1', '测试', 'test', '1', '测试', '2019-04-12 15:26:46', null);

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
INSERT INTO `sys_role_resource_rel` VALUES ('100', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('100', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('110', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('110', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('120', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('120', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('130', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('130', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('150', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('150', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('160', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('200', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('210', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('230', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('240', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('300', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('310', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('320', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('400', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('400', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('410', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('411', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('412', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('413', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('420', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('420', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('421', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('421', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('422', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('422', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('423', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('423', '2');
INSERT INTO `sys_role_resource_rel` VALUES ('430', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('431', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('432', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('433', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('440', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('441', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('443', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('444', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('445', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('450', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('460', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('461', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('462', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('463', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('464', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('465', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('470', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('471', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('472', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('473', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('474', '1');
INSERT INTO `sys_role_resource_rel` VALUES ('475', '1');

-- ----------------------------
-- Table structure for sys_user
-- ----------------------------
DROP TABLE IF EXISTS `sys_user`;
CREATE TABLE `sys_user` (
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
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='系统管理用户表';

-- ----------------------------
-- Records of sys_user
-- ----------------------------
INSERT INTO `sys_user` VALUES ('1', '1', 'admin', '超级管理员', 'oD2Ou3h126sv7bje58Z+fA==', '1', '186989878678', 'test@163.com', '2019-03-27 21:25:07', '2019-07-04 10:12:07', '测试');

-- ----------------------------
-- Table structure for sys_user_role_rel
-- ----------------------------
DROP TABLE IF EXISTS `sys_user_role_rel`;
CREATE TABLE `sys_user_role_rel` (
  `role_id` bigint(20) NOT NULL COMMENT '角色id',
  `sys_user_id` bigint(20) NOT NULL COMMENT '用户主键id',
  PRIMARY KEY (`role_id`,`sys_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='管理员与角色关联表';

-- ----------------------------
-- Records of sys_user_role_rel
-- ----------------------------
INSERT INTO `sys_user_role_rel` VALUES ('1', '1');
