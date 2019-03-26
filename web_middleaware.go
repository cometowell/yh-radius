package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

var urls = []string{"/login", "/statics/.+", "/favicon.ico"}

func PermCheck(c *gin.Context) {
	url := c.Request.URL.Path
	isOkay := isPassableUrl(url)
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
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": 999,
		"msg": "无权限访问!",
	})
}

func isPassableUrl(url string) bool {
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