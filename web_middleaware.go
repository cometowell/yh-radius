package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"time"
)

var urls = []string{"/login", "/statics/.+", "/favicon.ico"}

func PermCheck(c *gin.Context) {
	url := c.Request.URL.Path
	isOkay := checkUrl(url)
	if isOkay  {
		return
	}

	if 1==1 {
		c.Next()
		return
	}

	sessionVal, err := c.Cookie(SessionName)
	if err != nil {
		loginTimeOut(c)
		return
	}

	session := GlobalSessionManager.Provider.ReadSession(sessionVal)
	if session == nil {
		loginTimeOut(c)
		return
	}

	//权限校验
	c.Next()
}

func loginTimeOut(c *gin.Context) {
	c.Set("errMsg","login timeout")
	c.SetCookie(SessionName, "", -1, "/", "", false, true)
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func checkUrl(url string) bool {
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

func FormatDateTime(datatime time.Time) string {
	return datatime.Format(DateTimeFormat)
}