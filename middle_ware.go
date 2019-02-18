package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	ACCESS_ACCEPT_REPLY_MSG = "authenticate success"
)

// 验证用户名，密码
func UserVerify(ctx *Context)  error {
	// 验证用户名
	attr, ok := ctx.request.RadiusAttrStringKeyMap["User-Name"]
	if !ok {
		return errors.New("user's account number or password is incorrect")
	}

	userName := attr.AttrStringValue
	fmt.Println(userName)

	// 验证密码
	if ctx.request.isChap {
		fmt.Println("CHAP用户认证结果：", chap("111111", &ctx.request))
	} else {
		fmt.Println("PAP用户认证结果：", pap("111111", "111111", ctx.request))
	}

	return nil
}


// 验证MAC地址绑定
func MacAddrVerify(ctx *Context) {

}

// 设置通用认证响应属性
func AuthSetCommonResponseAttr(ctx *Context) {

}


func authReply(cxt *Context) {
	reply := RadiusPackage {
		Code:ACCESS_ACCEPT_CODE,
		Identifier: cxt.request.Identifier,
		Authenticator:[16]byte{},
	}

	replyMessage := RadiusAttr{
		AttrType: 18,
		AttrValue: []byte(ACCESS_ACCEPT_REPLY_MSG),
	}

	replyMessage.Length()
	reply.AddRadiusAttr(replyMessage)
	reply.PackageLength()

	// TODO secret
	secret := "111111"
	authReplyAuthenticator(cxt.request.Authenticator, &reply, secret)
	cxt.listener.WriteToUDP(reply.ToByte(), cxt.dst)
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