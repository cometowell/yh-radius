package database

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"go-rad/common"
	"go-rad/logger"
)

var DataBaseEngine *xorm.Engine

func InitDb()  {
	// 初始化数据库连接
	config := common.GetConfig()
	var err error
	DataBaseEngine, err = xorm.NewEngine(config["database.type"].(string), config["database.url"].(string))
	productStage := config["product.stage"].(string)

	if productStage != gin.ReleaseMode {
		DataBaseEngine.ShowSQL(true)
	}

	if err != nil {
		logger.Logger.Fatalf("连接数据库发生错误：%v", err)
	}
}
