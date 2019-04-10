package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
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

	router.POST("/product/add", addProduct)
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
}

func login(c *gin.Context) {
	var manager SysManager
	c.ShouldBindJSON(&manager)
	manager.Password = encrypt(manager.Password)
	ok, _ := engine.Get(&manager)

	if !ok {
		c.JSON(http.StatusOK, newErrorJsonResult("username or password is incorrect"))
		return
	}

	session := GlobalSessionManager.CreateSession(c)
	manager.Password = ""
	session.SetAttr("manager", manager)

	c.JSON(http.StatusOK, newSuccessJsonResult("success", session.SessionId()))
}

func logout(c *gin.Context) {
	GlobalSessionManager.DestroySession(c)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success"})
}

// -------------------------- manager start ----------------------------
func sessionManagerInfo(c *gin.Context) {
	token := c.GetHeader(SessionName)
	session := GlobalSessionManager.GetSession(token)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: session.GetAttr("manager")})
}

func managerList(c *gin.Context) {
	var params SysManager
	c.ShouldBindJSON(&params)
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
	c.ShouldBindJSON(&manager)
	engine.Id(manager.Id).Get(&manager)
	manager.Password = ""
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: manager})
}

func addManager(c *gin.Context) {
	var manager SysManager
	c.ShouldBindJSON(&manager)
	manager.Status = 1
	manager.Password = encrypt(manager.Password)
	manager.CreateTime = NowTime()
	engine.InsertOne(&manager)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "管理员信息添加成功!"})
}

func updateManager(c *gin.Context) {
	var manager SysManager
	c.ShouldBindJSON(&manager)
	if manager.Password != "" {
		manager.Password = encrypt(manager.Password)
	}
	engine.Id(manager.Id).Update(&manager)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "管理员信息更新成功!"})
}

func delManager(c *gin.Context) {
	var manager SysManager
	c.ShouldBindJSON(&manager)
	manager.Status = 3 // 标记为已删除
	engine.Id(manager.Id).Cols("status").Update(&manager)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "删除成功!"})
}

// -------------------------- manager end ----------------------------

// -------------------------- user start -----------------------------
func listUser(c *gin.Context) {
	var params RadUser
	c.ShouldBindJSON(&params)

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
	totalCount, err := engine.Table("rad_user").
		Alias("ru").Select(`ru.id,ru.username,ru.real_name,ru.product_id,
			ru.status,ru.available_time,ru.available_flow,ru.expire_time,
			ru.concurrent_count,ru.should_bind_mac_addr,ru.should_bind_vlan,ru.mac_addr,ru.vlan_id,
			ru.vlan_id2,ru.framed_ip_addr,ru.installed_addr,ru.mobile,ru.email,
			ru.pause_time,ru.create_time,ru.update_time,ru.description, sp.*`).
		Where(whereSql, whereArgs...).
		Limit(params.PageSize, (params.Page-1)*params.PageSize).
		Join("INNER", []string{"rad_product", "sp"}, "ru.product_id = sp.id").
		FindAndCount(&users)

	if err != nil {
		panic(err)
	}

	pagination := NewPagination(users, totalCount, params.Page, params.PageSize)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: pagination})
}

func fetchUserOrderRecord(c *gin.Context) {
	var user RadUser
	c.ShouldBindJSON(&user)
	var records []UserOrderRecordProduct
	engine.Join("INNER", "rad_product", "rad_product.id = user_order_record.product_id").Where("user_id = ?", user.Id).Asc("user_order_record.status").Find(&records)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: records})
}

func updateUser(c *gin.Context) {
	var user RadUser
	c.ShouldBindJSON(&user)
	session := engine.NewSession()
	defer session.Close()
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
	c.ShouldBindJSON(&user)
	user.Status = UserAvailableStatus
	user.CreateTime = NowTime()
	fmt.Printf("%#v", user)
	session := engine.NewSession()
	defer session.Close()
	var product RadProduct
	session.ID(user.ProductId).Get(&product)
	if product.Id == 0 {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"Code":    1,
			"Message": "产品不存在",
		})
		return
	}
	user.Password = encrypt(user.Password)
	purchaseProduct(&user, &product, c, session)
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
	c.ShouldBindJSON(&user)
	user.Status = UserDeletedStatus
	engine.Id(user.Id).Update(&user)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "删除成功!"})
}

// get user info by id
func fetchUser(c *gin.Context) {
	var user RadUser
	c.ShouldBindJSON(&user)
	engine.Id(user.Id).Get(&user)
	user.Password = ""
	c.JSON(http.StatusOK, defaultSuccessJsonResult(user))
}

func continueProduct(c *gin.Context) {
	var user RadUser
	c.ShouldBindJSON(&user)
	session := engine.NewSession()
	defer session.Close()

	bookOrderCount, e := session.Table("user_order_record").Where("user_id = ? and status = ?", user.Id, OrderBookStatus).Count()

	if e != nil {
		panic(e)
	}

	if bookOrderCount > 0 {
		c.JSON(http.StatusOK, newErrorJsonResult("用户已经预定了套餐暂未生效，不允许再次预定"))
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
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"Code":    1,
			"Message": "产品不存在",
		})
		return
	}
	if isExpire(oldUser.ExpireTime) { // 产品到期, 直接更新产品信息
		purchaseProduct(&oldUser, &newProduct, c, session)
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

func purchaseProduct(user *RadUser, product *RadProduct, c *gin.Context, session *xorm.Session) {
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

}

func updateProduct(c *gin.Context) {

}

func fetchProductList(c *gin.Context) {
	var products []RadProduct
	engine.Where("status = ?", 1).Find(&products)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: products})
}

func listProduct(c *gin.Context) {

}

func deleteProduct(c *gin.Context) {

}

func getProductInfo(c *gin.Context) {

}

// -------------------------- product end -----------------------------

// system
func fetchDepartments(c *gin.Context) {
	var departments []SysDepartment
	engine.Find(&departments)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: departments})
}
