package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	ACCESS_ACCEPT_REPLY_MSG = "authenticate success"
)

// 验证用户名，密码
func UserVerify(rp RadiusPackage)  {
	// 验证用户名

	// 验证密码
	if rp.isChap {
		fmt.Println("CHAP用户认证结果：", chap("111111", &rp))
	} else {
		fmt.Println("PAP用户认证结果：", pap("111111", "111111", rp))
	}
}


// 验证MAC地址绑定
func MacAddrVerify(rp RadiusPackage) {

}

// 设置通用认证响应属性
func AuthSetCommonResponseAttr(reply RadiusPackage) {

}


func authReply(rp RadiusPackage, listener *net.UDPConn, dest *net.UDPAddr) {
	reply := RadiusPackage {
		Code:ACCESS_ACCEPT_CODE,
		Identifier: rp.Identifier,
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
	authReplyAuthenticator(rp.Authenticator, &reply, secret)
	listener.WriteToUDP(reply.ToByte(), dest)
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