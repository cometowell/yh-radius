package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

var (
	engine  *xorm.Engine
	logger  *logrus.Logger
	config  map[string]interface{}
	limiter *rate.Limiter
)

type radEngine struct {
	radMiddleWares []RadMiddleWare
	port           int
	listener       *net.UDPConn
}

func Default(port int) (r *radEngine) {
	r = &radEngine{
		port: port,
	}
	r.radMiddleWares = append(r.radMiddleWares, RecoveryFunc(), NasValidation)
	return r
}

func (r *radEngine) Use(rms ...RadMiddleWare) {
	r.radMiddleWares = append(r.radMiddleWares, rms...)
}

func (r *radEngine) handlePackage(cxt context.Context) {

	UDPAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(r.port))
	if err != nil {
		logger.Fatalln("监听地址错误" + err.Error())
	}

	listener, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		logger.Fatalln("服务监听失败：", err)
	}

	r.listener = listener

	for {
		select {
		case <-cxt.Done():
			return
		default:
		}

		var pkg = make([]byte, MaxPackageLength)
		n, dst, err := listener.ReadFromUDP(pkg)
		if err != nil {
			logger.Infoln("接收认证请求报文发生错误", err.Error(), "消息来自 <<< ", dst.String())
			continue
		}

		if !limiter.Allow() {
			logger.Warn("服务器处理能力到底最高点：报文被丢弃", "消息来自 <<< ", dst.String())
			continue
		}

		go func(recPkg []byte, listener *net.UDPConn, dst *net.UDPAddr) {
			rp := parsePkg(recPkg)
			cxt := &Context{
				Request:  rp,
				Listener: listener,
				Dst:      dst,
				Handlers: r.radMiddleWares,
				index:    -1,
				Session:  engine.NewSession(),
			}

			cxt.Response = &RadiusPackage{
				Identifier:    cxt.Request.Identifier,
				Authenticator: [16]byte{},
			}

			//执行插件
			cxt.Next()
			if cxt.Response.Code != 0 {
				logger.Infof("响应报文：%+v\n", cxt.Response)
			}
		}(pkg[:n], listener, dst)
	}

}

func main() {
	bgCtx, cancelFunc := context.WithCancel(context.Background())
	// 加载配置文件
	config = loadConfig()
	logger = NewLogger()

	runtime.GOMAXPROCS(int(config["max.procs"].(float64)))

	// 处理协程数量限制
	limiter = rate.NewLimiter(rate.Limit(config["limiter.limit"].(float64)), int(config["limiter.burst"].(float64)))

	// 加载radius属性字典文件
	readAttributeFiles()
	logger.Info("字典文件加载完成, 正在启动radius服务...")

	// 初始化数据库连接
	var err error
	engine, err = xorm.NewEngine(config["database.type"].(string), config["database.url"].(string))
	engine.ShowSQL(true)
	if err != nil {
		logger.Fatalf("连接数据库发生错误：%v", err)
	}

	// 认证服务
	authServer := Default(int(config["authPort"].(float64)))
	authServer.Use(UserVerify)
	authServer.Use(VlanVerify)
	authServer.Use(MacAddrVerify)
	authServer.Use(AuthSpecAndCommonAttrSetter)
	authServer.Use(AuthAcceptReply)
	authServer.Use(TransactionCommitFunc)
	go authServer.handlePackage(bgCtx)
	logger.Info("已经启动Radius认证监听...")

	// 计费服务
	accountServer := Default(int(config["acctPort"].(float64)))
	accountServer.Use(AcctRecord)
	accountServer.Use(AcctReply)
	accountServer.Use(TransactionCommitFunc)
	go accountServer.handlePackage(bgCtx)
	logger.Info("已经启动Radius计费监听...")

	go webServer()
	logger.Info("已启动web服务")

	var scheduler = cron.New()
	err = scheduler.AddFunc(config["task.user.order.cron"].(string), userExpireTask)
	if err != nil {
		panic(err)
	}
	scheduler.Start()
	logger.Info("定时任务启动")

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
	logger.Warnln("RADIUS程序已退出...")
	os.Exit(0)
}

func webServer() {
	initWeb()
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(PermCheck)
	loadControllers(r)
	r.Run(config["web.server.url"].(string))
}

func loadConfig() map[string]interface{} {
	configBytes, err := ioutil.ReadFile("./config/radius.json")
	if err != nil {
		panic(err)
	}

	dst := make(map[string]interface{})
	json.Unmarshal(configBytes, &dst)
	return dst
}

func GetConfig() map[string]interface{} {
	return config
}
