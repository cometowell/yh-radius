package database

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"go-rad/common"
	"go-rad/logger"
)

var DataBaseEngine *xorm.Engine

func InitDb() {
	// 初始化数据库连接
	config := common.GetConfig()
	var err error
	DataBaseEngine, err = xorm.NewEngine(config["database.type"].(string), config["database.url"].(string))
	productStage := config["product.stage"].(string)

	if err != nil {
		logger.Logger.Fatalf("连接数据库发生错误：%v", err)
	}

	if productStage != gin.ReleaseMode {
		DataBaseEngine.ShowSQL(true)
	}
	// 连接失败则panic
	err = DataBaseEngine.Ping()
	if err != nil {
		logger.Logger.Fatalf("数据库连接失败,检查访问数据库用户名或密码是否正确: %s", err)
	}
}
