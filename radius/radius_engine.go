package radius

import (
	"context"
	"go-rad/common"
	"go-rad/database"
	"go-rad/logger"
	"golang.org/x/time/rate"
	"net"
	"strconv"
)

type radEngine struct {
	radMiddleWares []RadMiddleWare
	port           int
	listener       *net.UDPConn
	limiter *rate.Limiter
}

func Default(port int) (r *radEngine) {
	r = &radEngine{
		port: port,
	}
	r.radMiddleWares = append(r.radMiddleWares, RecoveryFunc(), NasValidation)
	// 处理协程数量限制
	r.limiter = rate.NewLimiter(rate.Limit(common.GetConfig()["limiter.limit"].(float64)), int(common.GetConfig()["limiter.burst"].(float64)))

	return r
}

func (r *radEngine) Use(rms ...RadMiddleWare) {
	r.radMiddleWares = append(r.radMiddleWares, rms...)
}

func (r *radEngine) HandlePackage(cxt context.Context) {

	UDPAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(r.port))
	if err != nil {
		logger.Logger.Fatalln("监听地址错误" + err.Error())
	}

	listener, err := net.ListenUDP("udp", UDPAddr)

	if err != nil {
		logger.Logger.Fatalln("服务监听失败：", err)
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
			logger.Logger.Infoln("接收认证请求报文发生错误", err.Error(), "消息来自 <<< ", dst.String())
			continue
		}

		if !r.limiter.Allow() {
			logger.Logger.Warn("服务器处理能力到底最高点：报文被丢弃", "消息来自 <<< ", dst.String())
			continue
		}

		go func(recPkg []byte, listener *net.UDPConn, dst *net.UDPAddr) {
			rp := parsePkg(recPkg)
			cxt := &Context{
				Request:  rp,
				Listener: listener,
				Dst:      dst,
				Handlers: r.radMiddleWares,
				Index:    -1,
				Session:  database.DataBaseEngine.NewSession(),
			}

			cxt.Response = &RadiusPackage{
				Identifier:    cxt.Request.Identifier,
				Authenticator: [16]byte{},
			}

			//执行插件
			cxt.Next()
			if cxt.Response.Code != 0 {
				logger.Logger.Infof("响应报文：%+v\n", cxt.Response)
			}
		}(pkg[:n], listener, dst)
	}

}
