package main

// radius codes constant,identifies the type of radius packet
const (

	// 认证请求
	ACCESS_REQUEST_CODE = 1

	// 认证响应
	ACCESS_ACCEPT_CODE = 2

	// 认证拒绝
	ACCESS_REJECT_CODE = 3

	// 计费请求
	ACCOUNTING_REQUEST_CODE = 4

	// 计费响应
	ACCOUNTING_RESPONSE_CODE = 5

	// 认证盘问
	ACCESS_CHALLENGE_CODE = 11

	// 往NAS发送下线通知
	DISCONNNECT_REQUEST = 40

	// NAS下线用户成功响应
	DISCONNECT_ACK = 41

	// NAS下线用户否定响应
	DISCONNECT_NAK = 42

	// COA授权改变请求
	COA_REQUEST = 43

	// 授权改变成功响应
	COA_ACK = 44

	// 授权改变否定响应
	COA_NAK = 45
)

const (

	// 认证字长度
	AUTHENTICATOR_LENGTH = 16

	// 最大radius报文长度
	MAX_PACKAGE_LENGTH = 4096

	// 26号私有属性Type
	VENDOR_SPECIFIC_TYPE byte = 26

	// 厂商ID长度
	VENDOR_ID_LENGTH = 4

	// radius属性头部长度
	// type(1) + length(1) = 2
	ATTR_HEADER_LENGHT = 2

	// 26号厂商私有属性头部长度 = 6
	// type(1) + length(1) + vendorId(4) = 6
	VENDOR_HEADER_LENGTH = ATTR_HEADER_LENGHT + VENDOR_ID_LENGTH

	// radius报文头部长度 = 20
	// code(1) + Identifier(1) + length(2) + Authenticator(16)
	PACKAGE_HEADER_LENGTH = 20
)