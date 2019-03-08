package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 注册Controller
func LoadController(router *gin.Engine) {

	// 用户管理
	router.GET("/user/list", UserList)
	router.GET("/user/add", UserAdd)
	router.GET("/user/insert", UserInsert)
	router.GET("/user/modify", UserModify)
	router.GET("/user/update", UserUpdate)

	router.GET("/login", LoginPage)
	//router.GET("/", LoginPage)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.GET("/index", Index)

}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}


func Login(c *gin.Context)  {
	var manager Manager
	err := c.ShouldBind(&manager)

	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	dbSession := engine.NewSession()
	dbSession.Begin()
	defer dbSession.Close()

	dbSession.Where("username = ?", manager.Username).Get(&manager)
	if manager.Id == 0 {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	session := GlobalSessionManager.CreateSession(c)
	fmt.Println(session)

	//加入菜单信息

	dbSession.Commit()
	c.Redirect(http.StatusMovedPermanently, "/index")
}

func Logout(c *gin.Context) {
	GlobalSessionManager.DestroySession(c)
	c.Redirect(http.StatusFound, "/login")
}

func UserAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "user_add.html", nil)
}

func UserInsert(c *gin.Context) {

}


func UserList(c *gin.Context) {
	c.HTML(http.StatusOK, "user_list.html", nil)
}

func UserModify(c *gin.Context) {

}

func UserUpdate(c *gin.Context) {

}