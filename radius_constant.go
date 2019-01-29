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
	VENDOR_SPECIFIC_TYPE byte = 26
	VENDOR_ID_LENGTH = 4

	// radius属性的type占字节长度
	ATTR_TYPE_FIELD_LENGHT = 1

	// radius属性的length占字节长度
	ATTR_LENGTH_FIELD_LENGHT = 1
)