package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
)

const (
	AccessAcceptReplyMsg = "authenticate success"
	ShouldBindMacAddrFlag = 1 // 用户绑定MAC地址标识
	ShouldBindVlanFlag 		  // 用户绑定虚拟专网标识
)

func NasValidation(cxt *Context) {
	nasIp := cxt.Dst.IP.String()
	logger.Infoln("UDP报文消息来源：", nasIp)
	logger.Infof("%+v\n", cxt.Request)

	nas := new(RadNas)
	engine.Where("ip_addr = ?", nasIp).Get(nas)
	// 验证UPD消息来源，非法来源丢弃
	if nas.Id == 0 {
		cxt.throwPackage = true
		panic("package come from unknown NAS: " + nasIp)
	}
	cxt.RadNas = *nas
	cxt.Next()
}

// 验证用户名，密码
func UserVerify(cxt *Context) {
	attr, ok := cxt.Request.RadiusAttrStringKeyMap["User-Name"]
	if !ok {
		userVerifyPanic()
	}

	userName := attr.AttrStringValue
	user := RadUser{UserName: userName}
	engine.Get(&user)

	if user.Id == 0 {
		userVerifyPanic()
	}

	// 验证密码
	password := decrypt(user.Password)
	if cxt.Request.isChap {
		if !chap(password, &cxt.Request) {
			userVerifyPanic()
		}
	} else {
		if !pap(cxt.RadNas.Secret, password, cxt.Request) {
			userVerifyPanic()
		}
	}
	cxt.User = &user
	cxt.Next()
}

func userVerifyPanic() {
	panic("user's account number or password is incorrect")
}

// 验证MAC地址绑定
func MacAddrVerify(cxt *Context) {
	if cxt.User.ShouldBindMacAddr == ShouldBindMacAddrFlag {
		attr, ok := cxt.Request.RadiusAttrStringKeyMap["Calling-Station-Id"]
		fmt.Println(attr, ok)
	}
	cxt.Next()
}

// 验证VLAN
func VlanVerify(cxt *Context) {
	if cxt.User.ShouldBindMacAddr == ShouldBindVlanFlag {
		attr, ok := cxt.Request.RadiusAttrStringKeyMap["NAS-Port-Id"]
		fmt.Println(attr, ok)
	}
	cxt.Next()
}

func getVlanIds(nasPortId string) (int, int) {

	ptn, _ := regexp.Compile(`vlanid=(\d);vlanid2=(\d)`)
	retMatch := ptn.FindStringSubmatch(nasPortId)
	vlanId := 0
	vlanId2 := 0
	var err error

	if len(retMatch) > 1 {
		vlanId, err = strconv.Atoi(retMatch[1])
	}

	if len(retMatch) > 2 {
		vlanId2, err = strconv.Atoi(retMatch[2])
	}

	if err != nil {
		return 0, 0
	}

	return vlanId, vlanId2
}

// 设置通用认证响应属性
func AuthSetCommonResponseAttr(cxt *Context) {
	// TODO 根据不同的厂商设置不同的响应属性
	cxt.Next()
}

func RecoveryFunc() RadMiddleWare {
	return func(cxt *Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorln("recovery invoke", err)

				if cxt.throwPackage {
					logger.Errorf("throw away package from %s: %+v\n", cxt.RadNas.IpAddr, cxt.Request)
					return
				}

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
