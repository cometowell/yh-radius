drop table if exists rad_user;
create table if not exists rad_user(
	id bigint(20) primary key comment '主键',
	username varchar(64) not null comment '账号',
  real_name varchar(128) comment '姓名',
	password varchar(256) not null comment '密码',
  product_id bigint(20) not null comment '产品ID',
	status int not null comment '状态，1：正常，2：停机，3：销户，4：禁用',
	available_time BIGINT(20) not null DEFAULT 0 comment '可用时长，单位：秒',
  available_flow BIGINT(20) not null DEFAULT 0 comment '可用流量，单位KB',
  expire_time datetime comment '到期时间',
	concurrent_count int not null DEFAULT 0 comment '并发数',
  should_bind_mac_addr tinyint(1) not null default 0 comment '需要绑定MAC地址，0：N，1：Y',
  should_bind_vlan tinyint(1) not null default 0 comment '需要绑定虚拟局域网，0：N，1：Y',
  mac_addr varchar(19) comment 'MAC地址',
  vlan_id int DEFAULT 0 comment '内层vlanId',
  vlan_id2 int DEFAULT 0  comment '外层vlanId',
	framed_ip_addr varchar(15) comment '用户绑定的静态IP地址',
	install_addr varchar(256) comment '装机地址',
  pause_time datetime comment '最近停机时间',
  create_time datetime not null comment '创建时间',
  update_time datetime comment '更新时间',
	description varchar(512) comment '描述'
) comment '用户表';

drop table if exists rad_user_wallet;
create table rad_user_wallet(
	id bigint(20) primary key,
  user_id bigint(20) UNIQUE KEY comment '用户ID',
	payment_password varchar(256) comment '支付密码'
) comment '用户钱包表';

drop table if exists rad_user_balance;
create table rad_user_balance(
	id bigint(20) primary key,
  user_wallet_id bigint(20) UNIQUE KEY comment '用户钱包ID',
	type int comment '类型1: 专项套餐，2：无限使用',
	product_id bigint comment '产品ID',
  balance int comment '余额',
  expire_time datetime comment '金额过期时间'
) comment '用户余额表';

drop table if exists rad_nas;
create table rad_nas(
	id bigint primary key,
  vendor_id int not null comment '厂商ID',
  name varchar(60) not null comment '名称',
  ip_addr varchar(15) not null comment 'IP地址',
	secret varchar(20) not null comment '共享秘钥',
	authorize_port int not null DEFAULT 3799 comment '授权端口，默认3799',
  description varchar(512) comment '描述'
) comment 'NAS网络接入设备表';

drop table if exists rad_product;
create table rad_product(
	id bigint primary key,
	name varchar(60) not null comment '产品名称',
  type int not null comment '产品类型,1：时长，2：流量',
  status int not null comment '状态,0:停用，1：正常',
	should_bind_mac_addr tinyint(1) not null default 0 comment '需要绑定MAC地址，0：N，1：Y',
  should_bind_vlan tinyint(1) not null default 0 comment '需要绑定虚拟局域网，0：N，1：Y',
	concurrent_count int not null DEFAULT 0 comment '并发数',
	product_duration bigint(20) not null default 0 comment '时长,单位秒',
  product_flow bigint(20) not null default 0 comment '流量，单位KB',
	flow_clear_cycle tinyint(1) not null comment '计费周期；0：无限时长， 1：日，2：月：3：固定（开通至使用时长截止[用户套餐过期时间]）',
	price int not null default 0 comment '产品价格，单位分',
  up_stream_limit bigint(20) not null comment '上行流量限制',
  down_stream_limit bigint(20) not null comment '下行流量限制',
	create_time datetime not null comment '创建时间',
	update_time datetime comment '创建时间',
  description varchar(512) comment '描述'
) comment '产品表';

drop table if exists online_user;
create table online_user(
	id bigint(2) primary key,
	username varchar(64) not null comment '账号',
	nas_ip_addr varchar(15) not null comment 'NAS IP地址',
	acct_session_id varchar(128) not null comment '计费session id',
	start_time datetime not null comment '开始计费时间',
	used_duration int not null default 0 comment '已计费时长',
	ip_addr varchar(15) not null comment '用户IP地址',
	mac_addr varchar(19) not null comment '用户MAC地址',
	NasPortId varchar(128) comment '标识用户认证端口',
	total_up_stream bigint(20) not null default 0 comment '上行总流量',
	total_down_stream bigint(20) not null default 0 comment '下行总流量'
) comment '在线用户表';