package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

const (
	ACCESS_ACCEPT_REPLY_MSG = "authenticate success"
)

func DefaultRecovery(cxt *Context) {

}

func NasValidation(cxt *Context) {
	nasIp := cxt.Dst.IP.String()
	fmt.Println(nasIp)
	cxt.Next()
}

// 验证用户名，密码
func UserVerify(cxt *Context) {
	panic("开始搞事情了")
	// 验证用户名
	attr, ok := cxt.Request.RadiusAttrStringKeyMap["User-Name"]
	if !ok {
		panic("user's account number or password is incorrect")
	}

	userName := attr.AttrStringValue
	fmt.Println(userName)

	// 验证密码
	if cxt.Request.isChap {

	} else {

	}
	cxt.Next()
}

// 验证MAC地址绑定
func MacAddrVerify(cxt *Context) {
	cxt.Next()
}

// 设置通用认证响应属性
func AuthSetCommonResponseAttr(cxt *Context) {
	cxt.Next()
}

func RecoveryFunc() RadMiddleWare {
	return func(cxt *Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovery invoke", err)
			}
		}()
		cxt.Next()
	}
}

func RecoveryWare() RadMiddleWare {
	return RecoveryFunc()
}

func AuthReply(cxt *Context) {
	cxt.Response = &RadiusPackage{
		Code:          AccessAcceptCode,
		Identifier:    cxt.Request.Identifier,
		Authenticator: [16]byte{},
	}

	replyMessage := RadiusAttr{
		AttrType:  18,
		AttrValue: []byte(ACCESS_ACCEPT_REPLY_MSG),
	}

	replyMessage.Length()
	cxt.Response.AddRadiusAttr(replyMessage)
	cxt.Response.PackageLength()

	// TODO secret
	secret := "111111"
	authReplyAuthenticator(cxt.Request.Authenticator, cxt.Response, secret)
	cxt.Listener.WriteToUDP(cxt.Response.ToByte(), cxt.Dst)
}

// ResponseAuth = MD5(Code+ID+Length+RequestAuth+Attributes+Secret)
func authReplyAuthenticator(authAuthenticator [16]byte, reply *RadiusPackage, secret string) {
	md5hash := md5.New()
	var buf bytes.Buffer
	buf.WriteByte(reply.Code)
	buf.WriteByte(reply.Identifier)

	var lengthBytes = make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, reply.Length)
	buf.Write(lengthBytes)

	buf.Write(authAuthenticator[:])

	for _, attr := range reply.RadiusAttrs {
		if attr.AttrLength == 0 {
			continue
		}
		buf.Write(attr.toBytes())
	}

	buf.Write([]byte(secret))

	md5hash.Write(buf.Bytes())
	sum := md5hash.Sum(nil)

	reply.Authenticator = getSixteenBytes(sum)
}
