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
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COMMENT='在线用户表';

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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='NAS网络接入设备表';

-- ----------------------------
-- Table structure for rad_product
-- ----------------------------
DROP TABLE IF EXISTS `rad_product`;
CREATE TABLE `rad_product` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(60) NOT NULL COMMENT '产品名称',
  `type` int(11) NOT NULL COMMENT '产品类型,1：时长，2：流量',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户上网记录';