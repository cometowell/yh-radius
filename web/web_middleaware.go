package web

import (
	"github.com/gin-gonic/gin"
	"go-rad/common"
	"go-rad/model"
	"net/http"
	"regexp"
)

var urls = []string{"/login", "/favicon.ico"}

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
	if isOkay {
		return
	}

	accessToken := c.GetHeader(common.SessionName)
	session := GlobalSessionManager.Provider.ReadSession(accessToken)
	if session == nil {
		notLoggedIn(c)
		return
	}

	c.Set("session", session)
	// 权限校验
	resources := session.GetAttr("resources").([]model.SysResource)
	canAccess := false

	url := c.Request.URL.Path
	for _, res := range resources {
		if res.Url == "" {
			continue
		}

		if res.Url == url {
			canAccess = true
			break
		}

		compile, _ := regexp.Compile(res.Url)
		if compile.Match([]byte(url)) {
			canAccess = true
			break
		}
	}

	if !canAccess {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    1,
			"message": "权限不足",
		})
		return
	}

	c.Next()
}

func notLoggedIn(c *gin.Context) {
	c.AbortWithStatus(http.StatusUnauthorized)
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
