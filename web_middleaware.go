package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

var urls = []string{"/login","logout", "/favicon.ico"}

func PermCheck(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	// options request
	httpMethod := c.Request.Method
	if httpMethod == http.MethodOptions {
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Context-Type", "application/json")
		c.Status(http.StatusOK)
		return
	}

	isOkay := isPassableUrl(c)
	if isOkay  {
		return
	}

	accessToken := c.GetHeader(SessionName)
	session := GlobalSessionManager.Provider.ReadSession(accessToken)
	if session == nil {
		unauthorizedAccess(c)
		return
	}

	// TODO 权限校验
	c.Next()
}

func unauthorizedAccess(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code": 999,
		"msg": "无权限访问!",
	})
}

func isPassableUrl(c *gin.Context) bool {
	url := c.Request.URL.Path
	for _, p := range urls {
		if p == url {
			return true
		}
		if ok, _ := regexp.MatchString(p, url); ok {
			return true
		}
	}
	return false
}