package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	AccessAcceptReplyMsg = "authenticate success"
)

func NasValidation(cxt *Context) {
	nasIp := cxt.Dst.IP.String()
	log.Println("UDP报文消息来源：", nasIp)
	log.Printf("%+v\n", cxt.Request)

	// 验证UPD消息来源，非法来源丢弃
	// cxt.throwPackage = true

	cxt.Next()
}

// 验证用户名，密码
func UserVerify(cxt *Context) {
	//panic("user's account number or password is incorrect")
	// 验证用户名
	attr, ok := cxt.Request.RadiusAttrStringKeyMap["User-Name"]
	if !ok {
		userVerifyPanic()
	}

	// TODO 验证用户名密码
	userName := attr.AttrStringValue
	fmt.Println(userName)

	// 验证密码
	if cxt.Request.isChap {
		if !chap("111111", &cxt.Request) {
			userVerifyPanic()
		}
	} else {
		if !pap("111111", "111111", cxt.Request) {
			userVerifyPanic()
		}
	}
	cxt.Next()
}

func userVerifyPanic() {
	panic("user's account number or password is incorrect")
}

// 验证MAC地址绑定
func MacAddrVerify(cxt *Context) {
	cxt.Next()
}

// 设置通用认证响应属性
func AuthSetCommonResponseAttr(cxt *Context) {
	// TODO 根据不同的厂商设置不同的响应属性
	cxt.Next()
}

func AuthRecoveryFunc() RadMiddleWare {
	return func(cxt *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("recovery invoke", err)
				if cxt.Request.Code == AccessRequestCode {
					authReply(cxt, AccessRejectCode, err.(string))
				}
			}
		}()
		cxt.Next()
	}
}

func AuthAcceptReply(cxt *Context) {
	authReply(cxt, AccessAcceptCode, AccessAcceptReplyMsg)
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
