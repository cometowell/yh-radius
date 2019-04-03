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

// ======================= manager start =========================
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
	pageByWhereSql(c, &managers, "" ,nil)
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
	manager.CreateTime = *NowTime()
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

// ======================= manager end =========================

// system
func fetchDepartments(c *gin.Context) {
	var departments []SysDepartment
	engine.Find(&departments)
	c.JSON(http.StatusOK, JsonResult{Code: 0, Message: "success", Data: departments})
}