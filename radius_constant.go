package main

// radius codes constant,identifies the type of radius packet
const (
	// 认证请求
	AccessRequestCode = 1

	// 认证响应
	AccessAcceptCode = 2

	// 认证拒绝
	AccessRejectCode = 3

	// 计费请求
	AccountingRequestCode = 4

	// 计费响应
	AccountingResponseCode = 5

	// 认证盘问
	AccessChallengeCode = 11

	// 往NAS发送下线通知
	DisconnectRequest = 40

	// NAS下线用户成功响应
	DisconnectAck = 41

	// NAS下线用户否定响应
	DisconnetNak = 42

	// COA授权改变请求
	COARequest = 43

	// 授权改变成功响应
	COAAck = 44

	// 授权改变否定响应
	COANak = 45

	CHAPPasswordType = 3
	CHAPChallenge    = 60
)

const (
	// 认证字长度
	AuthenticatorLength = 16

	// 最大radius报文长度
	MaxPackageLength = 4096

	// 26号私有属性Type
	VendorSpecificType byte = 26

	// 厂商ID长度
	VendorIdLength = 4

	// radius属性头部长度
	// type(1) + length(1) = 2
	AttrHeaderLength = 2

	// 26号厂商私有属性头部长度 = 6
	// type(1) + length(1) + vendorId(4) = 6
	VendorHeaderLength = AttrHeaderLength + VendorIdLength

	// radius报文头部长度 = 20
	// code(1) + Identifier(1) + length(2) + Authenticator(16)
	PackageHeaderLength = 20
)

const (
	Standard = 0
	Cisco    = 9
	Huawei   = 2011
	Zte      = 3902
	MikroTik = 14988
)

const (
	ClacPriceByMonth    = 1 // 按月收费
	UseTimesProductType = 2 // 自由时长
	UseFlowsProductType = 3 // 流量
)

const (
	AcctStatusTypeStart         = 1
	AcctStatusTypeStop          = 2
	AcctStatusTypeInterimUpdate = 3
	AcctStatusTypeAccountingOn  = 7
	AcctStatusTypeAccountingOff = 8
)

const DateTimeFormat = "2006-01-02 15:04:05"
