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
	log.Println("已经启动认证监听")
	for {
		var pkg= make([]byte, MAX_PACKAGE_LENGTH)
		n, sAddr, err := authListener.ReadFromUDP(pkg)
		if err != nil {
			log.Println("这里发生错误了", err.Error(), "消息来自 <<< ", sAddr.String())
		}

		// 这里需要控制协程的数量
		go handleAuth(pkg[:n])
	}

}

// 认证报文处理,认证 + 授权
func handleAuth(recvPkg []byte) {
	rp := parsePkg(recvPkg)
	fmt.Printf("%+v\n", rp)
}

func accountServer(accountListener *net.UDPConn) {
	log.Println("已经启动计费监听")
	for {
		// TODO 异步处理报文
		var pkg= make([]byte, MAX_PACKAGE_LENGTH)
		n, sAddr, err := accountListener.ReadFromUDP(pkg)
		fmt.Println("有UDP包来了")
		if err != nil {
			fmt.Println("这里发生错误了" + err.Error())
		}
		fmt.Println(n, sAddr.String(), err)
	}
}

// 计费报文处理
func handleAccounting() {

}