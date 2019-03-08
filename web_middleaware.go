package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

var urls  =  []string{"/login","/statics/.+","/favicon.ico"}

func PermCheck(c *gin.Context) {
	url := c.Request.URL.Path
	isOkay := checkUrl(url)
	if isOkay {
		return
	}

	sessionVal, err := c.Cookie(SessionName)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	session := GlobalSessionManager.Provider.ReadSession(sessionVal)
	if session == nil {
		c.SetCookie(SessionName, "",  -1, "/", "", false, true)
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	//权限校验


	c.Next()
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
