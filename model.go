package main

import (
	"time"
)

type RadUser struct {
	Id                int64  `xorm:"pk autoincr"`
	UserName          string `xorm:"unique 'username'"`
	RealName          string
	Password          string
	ProductId         int64
	Status            int
	AvailableTime     int   // sec
	AvailableFlow     int64 // KB
	ExpireTime        *time.Time
	ConcurrentCount   int // 并发数
	ShouldBindMacAddr int
	ShouldBindVlan    int
	MacAddr           string
	VlanId            int
	VlanId2           int
	FramedIpAddr      string // 静态IP
	InstalledAddr     string
	PauseTime         *time.Time // 停机时间
	CreateTime        *time.Time
	UpdateTime        *time.Time
	Description       string

	sessionTimeout int        `xorm:"-"`
	product        RadProduct `xorm:"-"`
}

type RadUserWallet struct {
	Id              int64 `xorm:"pk autoincr"`
	UserId          int64
	PaymentPassword string
	Balance         int
}

type RadUserSpecialBalance struct {
	Id           int64 `xorm:"pk autoincr"`
	UserWalletId int64
	Type         int // 1: 专项套餐，2：无限使用
	ProductId    int64
	Balance      int
	ExpireTime   time.Time
}

type OnlineUser struct {
	Id              int64  `xorm:"pk autoincr"`
	UserName        string `xorm:"'username'"`
	NasIpAddr       string
	AcctSessionId   string
	StartTime       time.Time
	UsedDuration    int //已记账时长:sec
	IpAddr          string
	MacAddr         string
	NasPortId       string // vlanid, vlanid2
	TotalUpStream   int64
	TotalDownStream int64
}

type RadProduct struct {
	Id                int64 `xorm:"pk autoincr"`
	Name              string
	Type              int // 类型：1:包月 2：自由时长，3：流量
	Status            int
	ShouldBindMacAddr int
	ShouldBindVlan    int
	ConcurrentCount   int
	ServiceMonth      int
	ProductDuration   int64 // 套餐使用时长：sec
	ProductFlow       int64 // 套餐流量 KB
	FlowClearCycle    int   // 流量清零周期；0：无限时长， 1：日，2：月：3：固定（开通至使用时长截止[用户套餐过期时间]）
	Price             int   //分
	UpStreamLimit     int   // 上行流量，Kb
	DownStreamLimit   int   // 下行流量，Kb
	DomainName        string
	Description       string
	CreateTime        *time.Time
	UpdateTime        *time.Time
}

type RadNas struct {
	Id            int64 `xorm:"pk autoincr"`
	VendorId      int
	Name          string
	IpAddr        string
	Secret        string
	AuthorizePort int //授权端口
	Description   string
}

type UserOnlineLog struct {
	Id              int64  `xorm:"pk autoincr"`
	UserName        string `xorm:"'username'"`
	StartTime       time.Time
	StopTime        *time.Time
	UsedDuration    int
	TotalUpStream   int
	TotalDownStream int
	NasIpAddr       string
	IpAddr          string
	MacAddr         string
}

type SysManager struct {
	Id           int64 `xorm:"pk autoincr"`
	DepartmentId int64
	Username     string
	Password     string
	RealName     string
	Status       int8
	Mobile       string
	Email        string
	CreateTime   time.Time
	UpdateTime   *time.Time
	Description  string
}

type SysDepartment struct {
	Id          int64 `xorm:"pk autoincr"`
	Code        string
	Name        string
	ParentId    int64
	CreateTime  time.Time
	UpdateTime  *time.Time
	Description string
}

type SysRole struct {
	Id          int64 `xorm:"pk autoincr"`
	Code        string
	Name        string
	CreateTime  time.Time
	UpdateTime  *time.Time
	Description string
}

type SysResource struct {
	Id                int64 `xorm:"pk autoincr"`
	ParentId          int64
	Name              string
	Icon              string
	Url               string
	Type              int
	Enable            int
	PermMark          string
	SortOrder         int
	Description       string
	ShouldPermControl int
}

type SysManagerRoleRel struct {
	ManagerId int64 `xorm:"pk"`
	RoleId    int64 `xorm:"pk"`
}

type SysRoleResourceRel struct {
	ResourceId int64 `xorm:"pk"`
	RoleId     int64 `xorm:"pk"`
}

type SysManagerRole struct {
	SysManager         `xorm:"extends"`
	SysManagerRoleRel  `xorm:"extends"`
	SysRole            `xorm:"extends"`
	SysRoleResourceRel `xorm:"extends"`
	SysResource        `xorm:"extends"`
}

func (SysManagerRole) TableName() string {
	return "sys_manager"
}
