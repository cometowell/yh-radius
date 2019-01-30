package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	authUDPAddr, err := net.ResolveUDPAddr("udp", ":1812")
	accountUDPAddr, err := net.ResolveUDPAddr("udp", ":1813")
	if err != nil {
		panic("监听地址错误" + err.Error())
	}

	// 认证监听
	authListener, err := net.ListenUDP("udp", authUDPAddr)
	accountListener, err := net.ListenUDP("udp", accountUDPAddr)

	if err != nil {
		panic("监听UDP连接失败,{}" + err.Error())
	}

	// 处理认证报文
	go handlerAuthPackage(authListener)

	// 处理计费报文
	go handlerAccountPackage(accountListener)

	// TODO 优雅关闭服务

	// 防止主线程退出
	select {}
}

func handlerAuthPackage(authListener *net.UDPConn) {
	fmt.Println("已经启动认证监听", time.Now())
	defer authListener.Close()

	for {

		// TODO 异步处理报文

		var pkg= make([]byte, MAX_PACKAGE_LENGTH)
		n, sAddr, err := authListener.ReadFromUDP(pkg)
		fmt.Println("有UDP包来了")
		if err != nil {
			fmt.Println("这里发生错误了" + err.Error())
		}

		fmt.Println(n, sAddr.String(), err)
	}

}

func handlerAccountPackage(accountListener *net.UDPConn) {
	fmt.Println("已经启动计费监听", time.Now())
	defer accountListener.Close()

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