package main

type RadUserProduct struct {
	RadUser    `xorm:"extends" json:"radUser"`
	RadProduct `xorm:"extends" json:"radProduct"`
}

func (RadUserProduct) TableName() string {
	return "rad_user"
}

type RadUser struct {
	Id                int64  `xorm:"pk autoincr" json:"id"`
	UserName          string `xorm:"unique 'username'" json:"username"`
	RealName          string `json:"realName"`
	Password          string `json:"password"`
	ProductId         int64  `json:"productId"`
	Status            int    `json:"status"`
	AvailableTime     int64  `json:"availableTime"` // sec
	AvailableFlow     int64  `json:"availableFlow"` // KB
	ExpireTime        Time   `json:"expireTime"`
	ConcurrentCount   int    `json:"concurrentCount"` // 并发数
	ShouldBindMacAddr int    `json:"shouldBindMacAddr"`
	ShouldBindVlan    int    `json:"shouldBindVlan"`
	MacAddr           string `json:"macAddr"`
	VlanId            int    `json:"vlanId"`
	VlanId2           int    `json:"vlanId2"`
	FramedIpAddr      string `json:"framedIpAddr"` // 静态IP
	InstalledAddr     string `json:"installedAddr"`
	PauseTime         Time   `json:"pauseTime"` // 停机时间
	CreateTime        Time   `json:"createTime"`
	UpdateTime        Time   `json:"updateTime"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
	Description       string `json:"description"`

	Product        RadProduct `xorm:"-" json:"product"`
	sessionTimeout int        `xorm:"-"`
	Pager          `xorm:"-" json:"page"`
	Count          int `xorm:"-" json:"count"`
	Price          int `xorm:"-" json:"price"`
}

type RadUserWallet struct {
	Id              int64  `xorm:"pk autoincr" json:"id"`
	UserId          int64  `json:"userId"`
	PaymentPassword string `json:"paymentPassword"`
	Balance         int    `json:"balance"`
}

type RadUserSpecialBalance struct {
	Id           int64 `xorm:"pk autoincr" json:"id"`
	UserWalletId int64 `json:"userWalletId"`
	Type         int   `json:"type"` // 1: 专项套餐，2：无限使用
	ProductId    int64 `json:"productId"`
	Balance      int   `json:"balance"`
	ExpireTime   Time  `json:"expireTime"`
}

type OnlineUser struct {
	Id              int64  `xorm:"pk autoincr" json:"id"`
	UserName        string `xorm:"'username'" json:"username"`
	NasIpAddr       string `json:"nasIpAddr"`
	AcctSessionId   string `json:"acctSessionId"`
	StartTime       Time   `json:"startTime"`
	UsedDuration    int    `json:"usedDuration"` //已记账时长:sec
	IpAddr          string `json:"ipAddr"`
	MacAddr         string `json:"macAddr"`
	NasPortId       string `json:"nasPortId"` // vlanid, vlanid2
	TotalUpStream   int64  `json:"totalUpStream"`
	TotalDownStream int64  `json:"totalDownStream"`
}

type RadProduct struct {
	Id                int64  `xorm:"pk autoincr" json:"id"`
	Name              string `json:"name"`
	Type              int    `json:"type"` // 类型：1:包月 2：自由时长，3：流量
	Status            int    `json:"status"`
	ShouldBindMacAddr int    `json:"shouldBindMacAddr"`
	ShouldBindVlan    int    `json:"shouldBindVlan"`
	ConcurrentCount   int    `json:"concurrentCount"`
	ServiceMonth      int    `json:"serviceMonth"`
	ProductDuration   int64  `json:"productDuration"` // 套餐使用时长：sec
	ProductFlow       int64  `json:"productFlow"`     // 套餐流量 KB
	FlowClearCycle    int    `json:"flowClearCycle"`  // 流量清零周期；0：无限时长， 1：日，2：月：3：固定（开通至使用时长截止[用户套餐过期时间]）
	Price             int    `json:"price"`           //分
	UpStreamLimit     int    `json:"upStreamLimit"`   // 上行流量，Kb
	DownStreamLimit   int    `json:"downStreamLimit"` // 下行流量，Kb
	DomainName        string `json:"domainName"`
	Description       string `json:"description"`
	CreateTime        Time   `json:"createTime"`
	UpdateTime        Time   `json:"updateTime"`

	Pager `xorm:"-" json:"page"`
}

type UserOrderRecord struct {
	Id        int64 `xorm:"pk autoincr" json:"id"`
	UserId    int64 `json:"userId"`
	ProductId int64 `json:"productId"`
	Price     int   `json:"price"`
	ManagerId int64 `json:"managerId"` // 操作人
	OrderTime Time  `json:"orderTime"`
	Status    int   `json:"status"`
	EndDate   Time  `json:"endDate"`
}

type UserOrderRecordProduct struct {
	UserOrderRecord `xorm:"extends" json:"userOrderRecord"`
	RadProduct      `xorm:"extends" json:"radProduct"`
}

func (UserOrderRecordProduct) TableName() string {
	return "user_order_record"
}

type RadNas struct {
	Id            int64  `xorm:"pk autoincr" json:"id"`
	VendorId      int    `json:"vendorId"`
	Name          string `json:"name"`
	IpAddr        string `json:"ipAddr"`
	Secret        string `json:"secret"`
	AuthorizePort int    `json:"authorizePort"` //授权端口
	Description   string `json:"description"`

	Pager `xorm:"-" json:"page"`
}

type UserOnlineLog struct {
	Id              int64  `xorm:"pk autoincr" json:"id"`
	UserName        string `xorm:"'username'" json:"username"`
	StartTime       Time   `json:"startTime"`
	StopTime        Time   `json:"stopTime"`
	UsedDuration    int    `json:"usedDuration"`
	TotalUpStream   int    `json:"totalUpStream"`
	TotalDownStream int    `json:"totalDownStream"`
	NasIpAddr       string `json:"nasIpAddr"`
	IpAddr          string `json:"ipAddr"`
	MacAddr         string `json:"macAddr"`
}

type SysManager struct {
	Id           int64  `xorm:"pk autoincr" json:"id"`
	DepartmentId int64  `json:"departmentId"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RealName     string `json:"realName"`
	Status       int8   `json:"status"`
	Mobile       string `json:"mobile"`
	Email        string `json:"email"`
	CreateTime   Time   `json:"createTime"`
	UpdateTime   Time   `json:"updateTime"`
	Description  string `json:"description"`

	Pager `xorm:"-" json:"page"`
}

type SysDepartment struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	ParentId    int64  `json:"parentId"`
	CreateTime  Time   `json:"createTime"`
	UpdateTime  Time   `json:"updateTime"`
	Description string `json:"description"`
}

type SysRole struct {
	Id          int64  `xorm:"pk autoincr" json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	CreateTime  Time   `json:"createTime"`
	UpdateTime  Time   `json:"updateTime"`
	Description string `json:"description"`
}

type SysResource struct {
	Id                int64  `xorm:"pk autoincr" json:"id"`
	ParentId          int64  `json:"parentId"`
	Name              string `json:"name"`
	Icon              string `json:"icon"`
	Url               string `json:"url"`
	Type              int    `json:"type"`
	Enable            int    `json:"enable"`
	PermMark          string `json:"permMark"`
	SortOrder         int    `json:"sortOrder"`
	Description       string `json:"description"`
	ShouldPermControl int    `json:"shouldPermControl"`
	Level             int    `json:"level"`

	Children []SysResource `xorm:"-" json:"children"`
}

type SysManagerRoleRel struct {
	ManagerId int64 `xorm:"pk" json:"managerId"`
	RoleId    int64 `xorm:"pk" json:"roleId"`
}

type SysRoleResourceRel struct {
	ResourceId int64 `xorm:"pk" json:"resourceId"`
	RoleId     int64 `xorm:"pk" json:"roleId"`
}

type SysManagerRole struct {
	SysManager         `xorm:"extends" json:"sysManager"`
	SysManagerRoleRel  `xorm:"extends" json:"sysManagerRoleRel"`
	SysRole            `xorm:"extends" json:"sysRole"`
	SysRoleResourceRel `xorm:"extends" json:"sysRoleResourceRel"`
	SysResource        `xorm:"extends" json:"sysResource"`
}

func (SysManagerRole) TableName() string {
	return "sys_manager"
}

const PageSize = 10

type Pagination struct {
	Size       int         `json:"size"`
	Current    int         `json:"current"`
	TotalPage  int64       `json:"totalPage"`
	TotalCount int64       `json:"totalCount"`
	Data       interface{} `json:"data"`
}

func NewPagination(data interface{}, totalCount int64, current, pageSize int) *Pagination {
	p := &Pagination{
		Size:       pageSize,
		Current:    current,
		Data:       data,
		TotalCount: totalCount,
	}
	p.setTotalPage()
	return p
}

func (p *Pagination) setTotalPage() {
	if p.TotalCount%int64(p.Size) != 0 {
		p.TotalPage = p.TotalCount/int64(p.Size) + 1
		return
	}
	p.TotalPage = p.TotalCount / int64(p.Size)
}

type Pager struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
