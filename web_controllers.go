package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func loadControllers(router *gin.Engine) {
	router.POST("/login", login)
	router.POST("/logout", logout)

	router.POST("/session/manager/info", sessionManagerInfo)
	router.POST("/manager/list", managerList)
	router.POST("/manager/info", managerById)
	router.POST("/manager/add", addManager)
	router.POST("/manager/update", updateManager)
	router.POST("/manager/delete", delManager)
	router.POST("/manager/change/password", changeManagerPassword)

	router.POST("/product/add", addProduct)
	router.POST("/product/list", listProduct)
	router.POST("/product/update", updateProduct)
	router.POST("/fetch/product", fetchProductList)
	router.POST("/product/info", getProductInfo)
	router.POST("/product/delete", deleteProduct)

	router.POST("/user/list", listUser)
	router.POST("/user/add", addUser)
	router.POST("/user/info", fetchUser)
	router.POST("/user/update", updateUser)
	router.POST("/user/order/record", fetchUserOrderRecord)
	router.POST("/user/continue", continueProduct)
	router.POST("/user/delete", deleteUser)

	router.POST("/fetch/department", fetchDepartments)
	router.POST("/department/info", getDepartmentInfo)
	router.POST("/department/list", listDepartments)
	router.POST("/department/add", addDepartment)
	router.POST("/department/update", updateDepartment)
	router.POST("/department/delete", deleteDepartment)

	router.POST("/nas/info", getNasInfo)
	router.POST("/nas/list", listNas)
	router.POST("/nas/add", addNas)
	router.POST("/nas/update", updateNas)
	router.POST("/nas/delete", deleteNas)

	router.POST("/resource/list", listRes)
	router.POST("/session/resource", getSessionResource)

	router.POST("/role/info", getRoleInfo)
	router.POST("/role/list", listRole)
	router.POST("/role/add", addRole)
	router.POST("/role/update", updateRole)
	router.POST("/role/delete", deleteRole)
	router.POST("/role/resources", empowerRole)
	router.POST("/role/empower/:roleId", doEmpowerRole)

	router.POST("/online/list", listOnline)
	router.POST("/online/off", offOnline)
	router.POST("/online/delete", deleteOnline)

}

func login(c *gin.Context) {
	var manager SysManager
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	manager.Password = encrypt(manager.Password)
	ok, _ := engine.Get(&manager)

	if !ok {
		c.JSON(http.StatusOK, newErrorJsonResult("username or password is incorrect"))
		return
	}

	session := GlobalSessionManager.CreateSession(c)
	manager.Password = ""
	session.SetAttr("manager", manager)

	var resources []SysResource
	engine.Table("sys_resource").Alias("sr").
		Join("LEFT", []string{"sys_role_resource_rel", "srr"}, "sr.id = srr.resource_id").
		Join("LEFT", []string{"sys_role", "r"}, "srr.role_id = r.id").
		Join("LEFT", []string{"sys_manager_role_rel", "smr"}, "smr.role_id = r.id").
		Join("LEFT", []string{"sys_manager", "m"}, "smr.manager_id = m.id").
		Where("m.id = ? or sr.should_perm_control = 0", manager.Id).
		Find(&resources)

	session.SetAttr("resources", resources)

	buttons := make([]int64, 0)
	for _, res := range resources {
		if res.Level == 3 {
			buttons = append(buttons, res.Id)
		}
	}

	c.JSON(http.StatusOK, newSuccessJsonResult("success", gin.H{
		"sessionId": session.SessionId(),
		"buttons":   buttons,
	}))
}

func logout(c *gin.Context) {
	GlobalSessionManager.DestroySession(c)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success"})
}

// -------------------------- manager start ----------------------------
func sessionManagerInfo(c *gin.Context) {
	token := c.GetHeader(SessionName)
	session := GlobalSessionManager.GetSession(token)

	managerInfo := session.GetAttr("manager")
	resources := session.GetAttr("resources").([]SysResource)

	buttons := make([]int64, 0)
	for _, res := range resources {
		if res.Level == 3 {
			buttons = append(buttons, res.Id)
		}
	}

	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: gin.H{
		"manager": managerInfo,
		"buttons": buttons,
	}})
}

func managerList(c *gin.Context) {
	var params SysManager
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.Set("current", params.Page)
	c.Set("pageSize", params.PageSize)
	var managers []SysManager
	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if params.Username != "" {
		whereSql += "and username like ? "
		whereArgs = append(whereArgs, "%"+params.Username+"%")
	}

	if params.RealName != "" {
		whereSql += "and real_name like ? "
		whereArgs = append(whereArgs, "%"+params.RealName+"%")
	}

	if params.Status != 0 {
		whereSql += "and status = ?"
		whereArgs = append(whereArgs, params.Status)
	}

	pageByWhereSql(c, &managers, whereSql, whereArgs)
}

func managerById(c *gin.Context) {
	var manager SysManager
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(manager.Id).Get(&manager)
	manager.Password = ""
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: manager})
}

func addManager(c *gin.Context) {
	var manager SysManager
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	count, _ := engine.Table("sys_manager").Where("username=?", manager.Username).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "用户名已存在!"})
		return
	}
	manager.Status = 1
	manager.Password = encrypt(manager.Password)
	manager.CreateTime = NowTime()
	engine.InsertOne(&manager)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "管理员信息添加成功!"})
}

func updateManager(c *gin.Context) {
	var manager SysManager
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	if manager.Password != "" {
		manager.Password = encrypt(manager.Password)
	}

	count, _ := engine.Table("sys_manager").Where("username=? and id != ?", manager.Username, manager.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "用户名已存在!"})
		return
	}
	manager.UpdateTime = NowTime()
	engine.Id(manager.Id).Update(&manager)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "管理员信息更新成功!"})
}

func delManager(c *gin.Context) {
	var manager SysManager
	err := c.ShouldBindJSON(&manager)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	manager.Status = 3 // 标记为已删除
	manager.UpdateTime = NowTime()
	engine.Id(manager.Id).Cols("status").Update(&manager)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "删除成功!"})
}

func changeManagerPassword(c *gin.Context) {
	var managerPassword ManagerPassword
	err := c.ShouldBindJSON(&managerPassword)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}

	var manager SysManager
	manager.Password = encrypt(managerPassword.NewPassword)
	_, err = engine.ID(managerPassword.Id).Cols("password").Update(&manager)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "密码修改成功!"})
}

// -------------------------- manager end ----------------------------

// -------------------------- user start -----------------------------
func listUser(c *gin.Context) {
	var params RadUser
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if params.UserName != "" {
		whereSql += "and ru.username like ? "
		whereArgs = append(whereArgs, "%"+params.UserName+"%")
	}

	if params.RealName != "" {
		whereSql += "and ru.real_name like ? "
		whereArgs = append(whereArgs, "%"+params.RealName+"%")
	}

	if params.Status != 0 {
		whereSql += "and ru.status = ?"
		whereArgs = append(whereArgs, params.Status)
	}

	var users []RadUserProduct
	totalCount, _ := engine.Table("rad_user").
		Alias("ru").Select(`ru.id,ru.username,ru.real_name,ru.product_id,
			ru.status,ru.available_time,ru.available_flow,ru.expire_time,
			ru.concurrent_count,ru.should_bind_mac_addr,ru.should_bind_vlan,ru.mac_addr,ru.vlan_id,
			ru.vlan_id2,ru.framed_ip_addr,ru.installed_addr,ru.mobile,ru.email,
			ru.pause_time,ru.create_time,ru.update_time,ru.description, sp.*`).
		Where(whereSql, whereArgs...).
		Limit(params.PageSize, (params.Page-1)*params.PageSize).
		Join("INNER", []string{"rad_product", "sp"}, "ru.product_id = sp.id").
		FindAndCount(&users)

	pagination := NewPagination(users, totalCount, params.Page, params.PageSize)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: pagination})
}

func fetchUserOrderRecord(c *gin.Context) {
	var user RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	var records []UserOrderRecordProduct
	engine.Join("INNER", "rad_product", "rad_product.id = user_order_record.product_id").Where("user_id = ?", user.Id).Asc("user_order_record.status").Find(&records)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: records})
}

func updateUser(c *gin.Context) {
	var user RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_user").Where("username = ? and id != ?", user.UserName, user.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名重复",
		})
		session.Rollback()
		return
	}

	var oldUser RadUser
	session.ID(user.Id).Get(&oldUser)
	// 停机用户重新使用需要顺延过期时间
	if oldUser.Status == UserPauseStatus && user.Status == UserAvailableStatus {
		hours := time.Now().Sub(time.Time(oldUser.PauseTime)).Hours()
		user.ExpireTime = Time(time.Time(user.ExpireTime).AddDate(0, 0, int(hours)/24))
	}
	session.ID(user.Id).Update(&user)
	session.Commit()
	c.JSON(http.StatusOK, newSuccessJsonResult("success", nil))
}

func addUser(c *gin.Context) {
	var user RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	user.Status = UserAvailableStatus
	user.CreateTime = NowTime()
	fmt.Printf("%#v", user)
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_user").Where("username = ?", user.UserName).Count()
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "用户名重复",
		})
		session.Rollback()
		return
	}

	var product RadProduct
	session.ID(user.ProductId).Get(&product)
	if product.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "产品不存在",
		})
		session.Rollback()
		return
	}
	user.Password = encrypt(user.Password)
	purchaseProduct(&user, &product)
	session.InsertOne(&user)
	// 订购信息
	webSession := GlobalSessionManager.GetSessionByGinContext(c)
	manager := webSession.GetAttr("manager").(SysManager)
	orderRecord := UserOrderRecord{
		UserId:    user.Id,
		ProductId: product.Id,
		Price:     user.Price,
		ManagerId: manager.Id,
		OrderTime: NowTime(),
		Status:    OrderUsingStatus,
		EndDate:   user.ExpireTime,
	}
	session.InsertOne(&orderRecord)

	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "用户添加成功!"})
}

func deleteUser(c *gin.Context) {
	var user RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	user.Status = UserDeletedStatus
	engine.Id(user.Id).Update(&user)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "删除成功!"})
}

// get user info by id
func fetchUser(c *gin.Context) {
	var user RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(user.Id).Get(&user)
	user.Password = ""
	c.JSON(http.StatusOK, defaultSuccessJsonResult(user))
}

func continueProduct(c *gin.Context) {
	var user RadUser
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	bookOrderCount, e := session.Table("user_order_record").Where("user_id = ? and status = ?", user.Id, OrderBookStatus).Count()

	if e != nil {
		session.Rollback()
	}

	if bookOrderCount > 0 {
		c.JSON(http.StatusOK, newErrorJsonResult("用户已经预定了套餐暂未生效，不允许再次预定"))
		session.Rollback()
		return
	}

	var oldUser RadUser
	session.ID(user.Id).Get(&oldUser)
	var newProduct RadProduct
	session.ID(user.ProductId).Get(&newProduct)

	var oldProduct RadProduct
	session.ID(oldUser.ProductId).Get(&oldProduct)

	webSession := GlobalSessionManager.GetSessionByGinContext(c)
	manager := webSession.GetAttr("manager").(SysManager)

	if newProduct.Id == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "产品不存在",
		})
		session.Rollback()
		return
	}
	if isExpire(oldUser.ExpireTime) { // 产品到期, 直接更新产品信息
		purchaseProduct(&oldUser, &newProduct)
	} else {
		// 产品未到期续订同一产品，修改过期时间
		if oldUser.ProductId == user.ProductId {
			oldUser.ExpireTime = Time(time.Time(oldUser.ExpireTime).AddDate(0, newProduct.ServiceMonth*user.Count, 0))
		} else {
			// 产品未到期续订不同产品，作为预定订单，当产品到期定时任务更换为预定产品
			expireTime, _ := getStdTimeFromString("2099-12-31 23:59:59")
			orderRecord := UserOrderRecord{
				UserId:    user.Id,
				ProductId: newProduct.Id,
				Price:     user.Price,
				ManagerId: manager.Id,
				OrderTime: NowTime(),
				Status:    OrderBookStatus,
				EndDate:   Time(expireTime),
			}
			session.InsertOne(&orderRecord)
		}
	}

	session.ID(oldUser.Id).Update(&oldUser)
	session.Commit()
}

func purchaseProduct(user *RadUser, product *RadProduct) {
	user.ShouldBindMacAddr = product.ShouldBindMacAddr
	user.ShouldBindVlan = product.ShouldBindVlan
	user.ConcurrentCount = product.ConcurrentCount
	user.AvailableTime = product.ProductDuration
	user.AvailableFlow = product.ProductFlow
	if product.Type == MonthlyProduct {
		expire := time.Time(user.ExpireTime)
		if time.Time(expire).IsZero() {
			expire = time.Now()
		}
		expire = time.Time(time.Date(expire.Year(), expire.Month()+time.Month(product.ServiceMonth), expire.Day(), 23, 59, 59, 0, expire.Location()))
		user.ExpireTime = Time(expire)
	} else if product.Type == TimeProduct {
		if time.Time(user.ExpireTime).IsZero() {
			expireTime, _ := getStdTimeFromString("2099-12-31 23:59:59")
			user.ExpireTime = Time(expireTime)
		}
	} else if product.Type == FlowProduct {
		if product.FlowClearCycle == DefaultFlowClearCycle {
			expireTime, _ := getStdTimeFromString("2099-12-31 23:59:59")
			user.ExpireTime = Time(expireTime)
		} else if product.FlowClearCycle == DayFlowClearCycle {
			user.ExpireTime = Time(getNextDayLastTime())
		} else if product.FlowClearCycle == MonthFlowClearCycle {
			user.ExpireTime = Time(getMonthLastTime())
		} else if product.FlowClearCycle == FixedPeriodFlowClearCycle {
			if time.Time(user.ExpireTime).IsZero() {
				user.ExpireTime = Time(getDayLastTimeAfterAYear())
			}
		}
	}
}

// -------------------------- user end ----------------------------------

// -------------------------- product start -----------------------------

func addProduct(c *gin.Context) {
	var product RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	product.Status = 1
	product.CreateTime = NowTime()
	engine.InsertOne(&product)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "套餐添加成功!"})
}

func updateProduct(c *gin.Context) {
	var product RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	product.UpdateTime = NowTime()
	session := engine.NewSession()
	session.Begin()
	engine.Id(product.Id).Update(&product)
	engine.Cols("concurrent_count").ID(product.Id).Update(&product)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "修改成功!"})
}

func fetchProductList(c *gin.Context) {
	var products []RadProduct
	engine.Where("status = ?", 1).Find(&products)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: products})
}

func listProduct(c *gin.Context) {
	var product RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if product.Name != "" {
		whereSql += "and name like ? "
		whereArgs = append(whereArgs, "%"+product.Name+"%")
	}

	if product.Status != 0 {
		whereSql += "and status = ? "
		whereArgs = append(whereArgs, product.Status)
	}

	if product.Type != 0 {
		whereSql += "and type = ?"
		whereArgs = append(whereArgs, product.Type)
	}

	var products []RadProduct
	totalCount, _ := engine.Table("rad_product").Where(whereSql, whereArgs...).
		Limit(product.PageSize, (product.Page-1)*product.PageSize).
		FindAndCount(&products)

	pagination := NewPagination(products, totalCount, product.Page, product.PageSize)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: pagination})
}

func deleteProduct(c *gin.Context) {
	var product RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	product.UpdateTime = NowTime()
	product.Status = 2
	engine.Cols("status, update_time").ID(product.Id).Update(&product)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "停用成功!"})
}

func getProductInfo(c *gin.Context) {
	var product RadProduct
	err := c.ShouldBindJSON(&product)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(product.Id).Get(&product)
	c.JSON(http.StatusOK, defaultSuccessJsonResult(product))
}

// -------------------------- product end -----------------------------

// -------------------------- nas start -----------------------------

func getNasInfo(c *gin.Context) {
	var nas RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(nas.Id).Get(&nas)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: nas})
}

func addNas(c *gin.Context) {
	var nas RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_nas").Where("ip_addr = ?", nas.IpAddr).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "错误：IP地址重复"})
		session.Rollback()
		return
	}
	session.InsertOne(&nas)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success"})
}

func updateNas(c *gin.Context) {
	var nas RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("rad_nas").Where("ip_addr = ? and id != ?", nas.IpAddr, nas.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "错误：IP地址重复"})
		session.Rollback()
		return
	}
	session.ID(nas.Id).Update(&nas)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "修改成功"})
}

func listNas(c *gin.Context) {
	var nas RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.Set("current", nas.Page)
	c.Set("pageSize", nas.PageSize)
	var nasList []RadNas
	pageByConditions(c, &nasList, &nas)
}

func deleteNas(c *gin.Context) {
	var nas RadNas
	err := c.ShouldBindJSON(&nas)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(nas.Id).Delete(&nas)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "已删除"})
}

// -------------------------- nas end -----------------------------

// -------------------------- resource start ---------------------------

func listRes(c *gin.Context) {
	var resList []SysResource
	engine.Find(&resList)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: getResLevel(resList)})
}

func getSessionResource(c *gin.Context) {
	c.JSON(http.StatusOK, defaultSuccessJsonResult(getManagerResources(c.GetHeader(SessionName))))
}

func getManagerResources(sessionId string) []SysResource {
	session := GlobalSessionManager.Provider.ReadSession(sessionId)
	resources := session.GetAttr("resources").([]SysResource)
	return getResLevel(resources)
}

// 菜单分层展示
func getResLevel(resList []SysResource) []SysResource {
	result := make([]SysResource, 0, 20)
	for _, res := range resList {
		if res.ParentId == 0 {
			r := res
			setChildren(&r, resList)
			result = append(result, r)
		}
	}
	return result
}

func setChildren(r *SysResource, resList []SysResource) {
	if r.Children == nil {
		r.Children = make([]SysResource, 0, 20)
	}
	for _, item := range resList {
		res := item
		if r.Id == res.ParentId {
			setChildren(&res, resList)
			r.Children = append(r.Children, res)
		}
	}
}

// -------------------------- resource end -----------------------------

// -------------------------- role start -------------------------------

func getRoleInfo(c *gin.Context) {
	var role SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(role.Id).Get(&role)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: role})
}

func addRole(c *gin.Context) {
	var role SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("sys_role").Where("code = ?", role.Code).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "错误：编码已存在，不能重复!"})
		session.Rollback()
		return
	}
	role.CreateTime = NowTime()
	role.Enable = 1
	session.InsertOne(&role)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "添加成功"})
}

func updateRole(c *gin.Context) {
	var role SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("sys_role").Where("code = ? and id != ?", role.Code, role.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "错误：编码重复!"})
		session.Rollback()
		return
	}
	session.ID(role.Id).Update(&role)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "修改成功"})
}

func listRole(c *gin.Context) {
	var role SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	c.Set("current", role.Page)
	c.Set("pageSize", role.PageSize)
	var roleList []SysRole

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if role.Name != "" {
		whereSql += "and name like ? "
		whereArgs = append(whereArgs, "%"+role.Name+"%")
	}

	if role.Code != "" {
		whereSql += "and code like ? "
		whereArgs = append(whereArgs, "%"+role.Code+"%")
	}
	pageByWhereSql(c, &roleList, whereSql, whereArgs)
}

func deleteRole(c *gin.Context) {
	var role SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(role.Id).Delete(&role)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "已删除"})
}

// 角色赋权
func empowerRole(c *gin.Context) {
	var role SysRole
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	var selectedResList []SysResource
	engine.Table("sys_resource").Alias("sr").
		Join("INNER", []string{"sys_role_resource_rel", "srr"}, "sr.id = srr.resource_id").
		Join("INNER", []string{"sys_role", "r"}, "srr.role_id = r.id").
		Where("r.id = ?", role.Id).
		Find(&selectedResList)

	var resources []SysResource
	engine.Where("should_perm_control = ?", 1).Find(&resources)
	for index, item := range resources {
		for _, r := range selectedResList {
			if r.Id == item.Id {
				resources[index].Selected = true
			}
		}
	}
	c.JSON(http.StatusOK, defaultSuccessJsonResult(getResLevel(resources)))
}

func doEmpowerRole(c *gin.Context) {
	roleId := c.Param("roleId")
	var roleResourceRels []SysRoleResourceRel
	err := c.ShouldBindJSON(&roleResourceRels)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	session.Where("role_id = ?", roleId).Delete(&SysRoleResourceRel{})
	if roleResourceRels != nil && len(roleResourceRels) > 0 {
		session.Insert(&roleResourceRels)
	}
	session.Commit()
	c.JSON(http.StatusOK, newSuccessJsonResult("赋权成功", nil))
}

// -------------------------- role end ---------------------------------

// department
func fetchDepartments(c *gin.Context) {
	var departments []SysDepartment
	engine.Find(&departments)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: departments})
}

func listDepartments(c *gin.Context) {
	var department SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if department.Code != "" {
		whereSql += "and sd.code like ? "
		whereArgs = append(whereArgs, "%"+department.Code+"%")
	}

	if department.Name != "" {
		whereSql += "and sd.name like ? "
		whereArgs = append(whereArgs, "%"+department.Name+"%")
	}

	if department.ParentId != 0 {
		whereSql += "and sd.parent_id = ? "
		whereArgs = append(whereArgs, department.ParentId)
	}

	var departments []Department
	count, _ := engine.Cols("sd.*, d.name").Table("sys_department").Alias("sd").
		Join("LEFT", []string{"sys_department", "d"}, "sd.parent_id = d.id").
		Where(whereSql, whereArgs...).
		Limit(department.PageSize, department.PageSize*(department.Page-1)).
		FindAndCount(&departments)

	pagination := NewPagination(departments, count, department.Page, department.PageSize)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: pagination})

}

func getDepartmentInfo(c *gin.Context) {
	var department SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	engine.Id(department.Id).Get(&department)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: department})
}

func addDepartment(c *gin.Context) {
	var department SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("sys_department").Where("code = ? or name = ?", department.Code, department.Name).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	department.Status = 1
	department.CreateTime = NowTime()
	session.InsertOne(&department)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "添加成功"})
}

func updateDepartment(c *gin.Context) {
	var department SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	count, _ := session.Table("sys_department").
		Where("(code = ? or name = ?) and id != ?", department.Code, department.Name, department.Id).Count()
	if count > 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "错误：编码或者名称重复"})
		session.Rollback()
		return
	}
	session.AllCols().ID(department.Id).Update(&department)
	session.Commit()
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "修改成功"})
}

func deleteDepartment(c *gin.Context) {
	var department SysDepartment
	err := c.ShouldBindJSON(&department)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error()})
		return
	}
	department.Status = 2 // 标记为停用
	engine.Id(department.Id).Cols("status").Update(&department)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "已停用!"})
}

// -------------------------- online --------------------------------
func listOnline(c *gin.Context) {
	var online OnlineUser
	c.ShouldBindJSON(&online)

	whereSql := "1=1 "
	whereArgs := make([]interface{}, 0)
	if online.UserName != "" {
		whereSql += "and ol.username like ? "
		whereArgs = append(whereArgs, "%"+online.UserName+"%")
	}

	if online.IpAddr != "" {
		whereSql += "and ol.ip_addr = ? "
		whereArgs = append(whereArgs, online.IpAddr)
	}

	if online.RealName != "" {
		whereSql += "and ru.real_name like ? "
		whereArgs = append(whereArgs, "%"+online.RealName+"%")
	}

	var onlines []Online
	count, _ := engine.Table("online_user").Alias("ol").
		Join("INNER", []string{"rad_user", "ru"}, "ol.username = ru.username").
		Where(whereSql, whereArgs...).
		Limit(online.PageSize, online.PageSize*(online.Page-1)).
		FindAndCount(&onlines)

	pagination := NewPagination(onlines, count, online.Page, online.PageSize)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: pagination})
}

func offOnline(c *gin.Context) {
	var online OnlineUser
	c.ShouldBindJSON(&online)

	var dst OnlineUser
	engine.ID(online.Id).Get(&dst)

	if dst.Id == 0 {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: "no records were found", Data: nil})
		return
	}

	//下线用户
	err := offlineUser(dst)
	if err != nil {
		c.JSON(http.StatusOK, JsonResult{Code: 1, Message: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, newSuccessJsonResult("下线成功", nil))
}

func deleteOnline(c *gin.Context) {
	var online OnlineUser
	c.ShouldBindJSON(&online)
	count, e := engine.ID(online.Id).Delete(&online)
	if e != nil || count == 0 {
		c.JSON(http.StatusOK, newErrorJsonResult("清理在线用户失败"))
		return
	}
	c.JSON(http.StatusOK, newSuccessJsonResult("清理在线用户成功", nil))
}
