package main

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func NasValidation(cxt *Context) {
	nasIp := cxt.Dst.IP.String()
	logger.Infoln("UDP报文消息来源：", nasIp)
	logger.Infof("%v\n", cxt.Request)
	cxt.Session.Begin()
	nas := new(RadNas)
	cxt.Session.Where("ip_addr = ?", nasIp).Get(nas)
	// 验证UPD消息来源，非法来源丢弃
	if nas.Id == 0 {
		cxt.throwPackage = true
		panic("package come from unknown NAS: " + nasIp)
	}

	// 验证
	cxt.RadNas = *nas
	cxt.Next()
}

func RecoveryFunc() RadMiddleWare {
	return func(cxt *Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorln("recovery invoke", err)
				if !cxt.Session.IsClosed() {
					cxt.Session.Rollback()
					cxt.Session.Close()
				}
				if cxt.throwPackage {
					logger.Errorf("throw away package from %s: %+v\n", cxt.RadNas.IpAddr, cxt.Request)
					if _, ok := err.(string); !ok {
						logger.Debug("异常堆栈信息：" + string(debug.Stack()))
					}
					return
				}

				if cxt.Request.Code == AccessRequestCode {
					var errMsg string
					if entry, ok := err.(*logrus.Entry); ok {
						errMsg = entry.Message
					} else if msg, ok := err.(string); ok {
						errMsg = msg
					} else {
						errMsg = "occur unknown error"
						logger.Errorf("occur unknown error: %+v", err)
						logger.Debug("异常堆栈信息：" + string(debug.Stack()))
					}
					authReply(cxt, AccessRejectCode, errMsg)
				}
			}
		}()
		cxt.Next()
	}
}

func TransactionCommitFunc(cxt *Context) {
	if cxt.Session != nil && !cxt.Session.IsClosed() {
		cxt.Session.Commit()
		cxt.Session.Close()
	}
}