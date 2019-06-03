package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"go-rad/common"
	"go-rad/database"
	"go-rad/logger"
	"go-rad/radius"
	"go-rad/task"
	"go-rad/web"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

func main() {
	bgCtx, cancelFunc := context.WithCancel(context.Background())
	// 加载配置文件
	config := common.GetConfig()
	runtime.GOMAXPROCS(int(config["max.procs"].(float64)))

	database.InitDb()

	// 加载radius属性字典文件
	radius.ReadAttributeFiles()
	logger.Logger.Info("字典文件加载完成, 正在启动radius服务...")

	// 认证服务
	authServer := radius.Default(int(config["auth.port"].(float64)))
	authServer.Use(radius.UserVerify)
	authServer.Use(radius.VlanVerify)
	authServer.Use(radius.MacAddrVerify)
	authServer.Use(radius.AuthSpecAndCommonAttrSetter)
	authServer.Use(radius.AuthAcceptReply)
	authServer.Use(radius.TransactionCommitFunc)
	go authServer.HandlePackage(bgCtx)
	logger.Logger.Info("已经启动Radius认证监听...")

	// 计费服务
	accountServer := radius.Default(int(config["acct.port"].(float64)))
	accountServer.Use(radius.AcctRecord)
	accountServer.Use(radius.AcctReply)
	accountServer.Use(radius.TransactionCommitFunc)
	go accountServer.HandlePackage(bgCtx)
	logger.Logger.Info("已经启动Radius计费监听...")

	go web.WebServer()
	logger.Logger.Info("已启动web服务")

	var scheduler = cron.New()
	err := scheduler.AddFunc(config["task.user.order.cron"].(string), task.UserExpireTask)
	if err != nil {
		panic(err)
	}
	scheduler.Start()
	logger.Logger.Info("定时任务启动")

	pid := syscall.Getpid()
	pidStr := strconv.Itoa(pid)
	err = ioutil.WriteFile("rad.pid", []byte(pidStr), 0777)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	cancelFunc()
	logger.Logger.Warnln("RADIUS程序已退出...")
	os.Exit(0)
}


