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
	router.GET("/", LoginPage)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	router.GET("/index", Index)

}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func LoginPage(c *gin.Context) {
	errMsg := c.Query("errMsg")
	data := make(map[string]interface{})
	if errMsg != "" {
		data["errMsg"] = errMsg
	}
	c.HTML(http.StatusOK, "login.html", data)
}

func Login(c *gin.Context) {
	if 1==1 {
		c.Redirect(http.StatusMovedPermanently, "/index")
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	dbSession := engine.NewSession()
	dbSession.Begin()
	defer dbSession.Close()

	manager := SysManager{Username: username}
	dbSession.Get(&manager)
	if manager.Id == 0 {
		loginError(c)
		return
	}

	if manager.Password != encrypt(password) {
		loginError(c)
		return
	}

	//加入权限信息

	session := GlobalSessionManager.CreateSession(c)
	fmt.Println(session)

	dbSession.Commit()
	c.Redirect(http.StatusMovedPermanently, "/index")
}

func loginError(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login" + buildUrlParams("errMsg", "username or password is incorrect"))
}

func Logout(c *gin.Context) {
	GlobalSessionManager.DestroySession(c)
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func UserAdd(c *gin.Context) {
	c.HTML(http.StatusOK, "user_add.html", nil)
}

func UserInsert(c *gin.Context) {

}

func UserList(c *gin.Context) {
	var userList []UserProduct
	engine.Table("rad_user").Alias("ru").
		Join("INNER", []string{"rad_product", "rp"}, "ru.product_id = rp.id").Find(&userList)
	c.HTML(http.StatusOK, "user_list.html", userList)
}

func UserModify(c *gin.Context) {

}

func UserUpdate(c *gin.Context) {

}
