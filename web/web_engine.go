package web

import (
	"container/list"
	"github.com/gin-gonic/gin"
	"go-rad/common"
)

func WebServer() {
	config := common.GetConfig()
	initWeb(config)
	gin.SetMode(config["product.stage"].(string))
	r := gin.Default()
	r.Use(PermCheck)
	loadControllers(r)
	r.Run(config["web.server.url"].(string))
}

func initWeb(config map[string]interface{}) {
	provider := &MemoryProvider{SesList: list.New()}
	providers["memory"] = provider
	provider.Sessions = make(map[string]*list.Element)
	gsm, err := CreateSessionManager(common.SessionName, "memory", int64(config["web.session.timeout"].(float64)))
	if err != nil {
		panic(err)
	}
	GlobalSessionManager = gsm
	go GlobalSessionManager.Gc()
}