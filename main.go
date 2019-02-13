package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// 读取radius属性字典文件
	readAttributeFiles()
	log.Println("字典文件加载完成...\n正在启动radius服务")
	// 启动radius服务
	server()

}

func server() {

	authUDPAddr, err := net.ResolveUDPAddr("udp", ":1812")
	accountUDPAddr, err := net.ResolveUDPAddr("udp", ":1813")
	if err != nil {
		panic("监听地址错误" + err.Error())
	}

	// 认证监听
	authListener, authErr := net.ListenUDP("udp", authUDPAddr)
	// 计费监听
	accountListener, accountErr := net.ListenUDP("udp", accountUDPAddr)

	if authErr != nil || accountErr != nil {
		log.Fatalln("认证服务或者计费服务监听失败：", authErr, accountErr)
	}

	defer authListener.Close()
	defer accountListener.Close()

	// 处理认证报文服务
	go authServer(authListener)

	// 处理计费报文服务
	go accountServer(accountListener)

	// TODO 优雅关闭服务

	// 防止主线程退出,监听退出信号
	select {}
}

func authServer(authListener *net.UDPConn) {
	log.Println("已经启动认证监听...")
	for {
		var pkg= make([]byte, MAX_PACKAGE_LENGTH)
		n, sAddr, err := authListener.ReadFromUDP(pkg)
		if err != nil {
			log.Println("接收认证请求报文发生错误", err.Error(), "消息来自 <<< ", sAddr.String())
			continue
		}

		// 这里需要控制协程的数量
		go handleAuth(pkg[:n], authListener, sAddr)
	}

}

// 认证报文处理,认证 + 授权
func handleAuth(recvPkg []byte, authListener *net.UDPConn, dest *net.UDPAddr) {
	rp := parsePkg(recvPkg)
	log.Printf("%+v\n", rp)

	// 认证用户信息, 中间件的形式处理
	for _, attr := range rp.RadiusAttrs {
		fmt.Println(attr)
	}

	//返回认证授权结果
	authReply(rp, authListener, dest)
}

func accountServer(accountListener *net.UDPConn) {
	log.Println("已经启动计费监听...")
	for {
		var pkg= make([]byte, MAX_PACKAGE_LENGTH)
		n, sAddr, err := accountListener.ReadFromUDP(pkg)
		if err != nil {
			fmt.Println("接收计费请求报文发送错误：", err.Error(), "消息来自 <<< ", sAddr.String() )
			continue
		}

		// 这里需要控制协程的数量
		go handleAccounting(pkg[:n], accountListener)
	}
}

// 计费报文处理
func handleAccounting(recvPkg []byte, accountListener *net.UDPConn) {
	rp := parsePkg(recvPkg)
	log.Printf("%+v\n", rp)
}