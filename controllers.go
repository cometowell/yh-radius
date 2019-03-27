package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func loadControllers(router *gin.Engine) {
	router.POST("/login", login)
	router.POST("/logout", logout)
}

func login(c *gin.Context) {
	rawData, _ := c.GetRawData()
	fmt.Println(string(rawData))
	c.JSON(http.StatusOK, gin.H{
		"ok": "bbb",
	})
}

func logout(c *gin.Context) {
	GlobalSessionManager.DestroySession(c)
	c.JSON(http.StatusOK, JsonResult{Code:1, Message:"你大爷的"})
}