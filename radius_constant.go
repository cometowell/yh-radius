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
)

const (

	MAX_PACKAGE_LENGTH = 4096

	// 26号私有属性Type
	VENDOR_SPECIFIC_TYPE byte = 26

	// 厂商ID长度
	VENDOR_ID_LENGTH = 4

	// type(1) + length(1) = 2
	ATTR_HEADER_LENGHT = 2

	// code(1) + Identifier(1) + length(2) + Authenticator(16)
	PACKAGE_HEADER_LENGTH = 20
)