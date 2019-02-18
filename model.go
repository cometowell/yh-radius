package main

import (
	"net"
	"time"
)

type Context struct {
	request RadiusPackage
	response *RadiusPackage
	User *User
	listener *net.UDPConn
	dst *net.UDPAddr
}

type User struct {
	Id uint64
	UserName string
	RealName string
	Password string
	ProductId uint64
	Status int
	RemainingTime uint64 // sec
	RemainingFlow uint64 // KB
	ExpireDate *time.Time
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
}

type UserWallet struct {
	Id uint64
	UserId uint64
	PaymentPassword string
	Balance uint
}

type UserSpecialBalance struct {
	UserWalletId uint64
	Type int // 1: 专项套餐，2：无限使用
	ProductId uint64
	Balance uint
	Expire time.Time
}

type OnlineUser struct {
	Id uint64
	UserName string
	NasIpAddr string
	SessionId string
	StartTime time.Time
	UsedDuration int //已记账时长:sec
	IpAddr string
	MacAddr string
	NasPortId string
	TotalUpStream uint64
	TotalDownStream uint64
}

type Product struct {
	Id uint64
	Name string
	Type int // 类型：0：时长，1：流量
	Status int
	ShouldBindMac int
	ShouldBindVlan int
	ConcurrentCount uint
	PackageDuration uint64 // 套餐使用时长：sec
	PackageFlow uint64 // 套餐流量 KB
	FlowClearCycle uint // 计费周期；0：无限时长， 1：日，2：月：3：固定（开通至使用时长截止[用户套餐过期时间]）
	Price uint //分
	UpStreamLimit uint64 // 上行流量，KB
	DownStreamLimit uint64 // 下行流量，KB
	Description string
	CreateTime *time.Time
	UpdateTime *time.Time
}

type Nas struct {
	Id uint64
	VendorId int
	Name string
	IpAddr string
	Secret string
	AuthorizePort int //授权端口
	Description string
}