package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func loadControllers(router *gin.Engine) {
	router.POST("/login", login)
	router.POST("/logout", logout)

	router.POST("/session/manager/info", sessionManagerInfo)
	router.POST("/manager/list", managerList)
	router.POST("/manager/info", managerById)
	router.POST("/manager/add", addManager)
	router.POST("/manager/update", updateManager)
	router.POST("/manager/del", delManager)

	router.POST("/user/list", listUser)

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

func updateUser(c *gin.Context) {

}

func addUser(c *gin.Context) {

}

func deleteUser(c *gin.Context) {

}

// -------------------------- user end -----------------------------

// system
func fetchDepartments(c *gin.Context) {
	var departments []SysDepartment
	engine.Find(&departments)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: departments})
}