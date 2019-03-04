package main

import (
	"time"
)

type RadUser struct {
	Id uint64 `xorm:"pk autoincr"`
	UserName string `xorm:"'username'"`
	RealName string
	Password string
	ProductId uint64
	Status int
	AvailableTime uint64 // sec
	AvailableFlow uint64 // KB
	ExpireTime *time.Time
	ConcurrentCount uint // 并发数
	ShouldBindMacAddr int
	ShouldBindVlan int
	MacAddr string
	VlanId int
	VlanId2 int
	FramedIpAddr string // 静态IP
	InstalledAddr string
	PauseTime *time.Time // 停机时间
	CreateTime *time.Time
	UpdateTime *time.Time
	Description string
}

type RadUserWallet struct {
	Id uint64 `xorm:"pk autoincr"`
	UserId uint64
	PaymentPassword string
	Balance uint
}

type RadUserSpecialBalance struct {
	Id uint64 `xorm:"pk autoincr"`
	UserWalletId uint64
	Type int // 1: 专项套餐，2：无限使用
	ProductId uint64
	Balance uint
	ExpireTime time.Time
}

type OnlineUser struct {
	Id uint64 `xorm:"pk autoincr"`
	UserName string
	NasIpAddr string
	AccSessionId string
	StartTime time.Time
	UsedDuration int //已记账时长:sec
	IpAddr string
	MacAddr string
	NasPortId string // vlanid, vlanid2
	TotalUpStream uint64
	TotalDownStream uint64
}

type RadProduct struct {
	Id uint64 `xorm:"pk autoincr"`
	Name string
	Type int // 类型：1：时长，2：流量
	Status int
	ShouldBindMacAddr int
	ShouldBindVlan int
	ConcurrentCount uint
	ProductDuration uint64 // 套餐使用时长：sec
	ProductFlow uint64 // 套餐流量 KB
	FlowClearCycle uint // 计费周期；0：无限时长， 1：日，2：月：3：固定（开通至使用时长截止[用户套餐过期时间]）
	Price uint //分
	UpStreamLimit uint32 // 上行流量，KB
	DownStreamLimit uint32 // 下行流量，KB
	DomainName string
	Description string
	CreateTime *time.Time
	UpdateTime *time.Time
}

type RadNas struct {
	Id uint64 `xorm:"pk autoincr"`
	VendorId int
	Name string
	IpAddr string
	Secret string
	AuthorizePort int //授权端口
	Description string
}