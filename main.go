package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"log"
	"net"
	"strconv"
)

var engine *xorm.Engine

type radEngine struct {
	radMiddleWares []RadMiddleWare
	port int
	listener *net.UDPConn
}

func Default(port int) (r *radEngine) {
	r = &radEngine{
		port: port,
	}
	r.radMiddleWares = append(r.radMiddleWares, NasValidation)
	return r
}

func (r *radEngine) Use(rmw RadMiddleWare) {
	r.radMiddleWares = append(r.radMiddleWares, rmw)
}

func (r *radEngine) handlePackage() {

	UDPAddr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(r.port))
	if err != nil {
		log.Fatalln("监听地址错误" + err.Error())
	}

	listener, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		log.Fatalln("服务监听失败：", err)
	}

	r.listener = listener

	for {
		var pkg = make([]byte, MaxPackageLength)
		n, dst, err := listener.ReadFromUDP(pkg)
		if err != nil {
			log.Println("接收认证请求报文发生错误", err.Error(), "消息来自 <<< ", dst.String())
			continue
		}

		// 这里需要控制协程的数量
		go func(recPkg []byte, listener *net.UDPConn, dst *net.UDPAddr) {
			rp := parsePkg(recPkg)
			cxt := &Context {
				Request: rp,
				Listener: listener,
				Dst: dst,
				Handlers: r.radMiddleWares,
				index: -1,
			}
			//执行插件
			cxt.Next()

		} (pkg[:n], listener, dst)
	}

}

func main() {

	// 加载配置文件
	config := loadConfig()

	// 加载radius属性字典文件
	readAttributeFiles()
	log.Println("字典文件加载完成, 正在启动radius服务...")

	// 初始化数据库连接
	var err error
	engine, err = xorm.NewEngine("mysql", "root:root@/127.0.0.1:3306/radius?charset=utf8")
	if err != nil {
		log.Fatalf("连接数据库发生错误：%v", err)
	}

	// 认证服务
	authServer := Default(int(config["authPort"].(float64)))
	authServer.Use(AuthRecoveryFunc())
	authServer.Use(UserVerify)
	authServer.Use(AuthAcceptReply)
	go authServer.handlePackage()
	log.Println("已经启动Radius认证监听...")

	// 计费服务
	accountServer := Default(int(config["acctPort"].(float64)))
	accountServer.Use(nil)
	go accountServer.handlePackage()
	log.Println("已经启动Radius计费监听...")

	// TODO 优雅关闭服务

	// 防止主线程退出,监听退出信号
	select {}
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
