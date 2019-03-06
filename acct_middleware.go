package main

import (
	"encoding/binary"
	"strconv"
	"time"
)

func AcctReply(cxt *Context) {
	cxt.Response.Code = AccountingResponseCode
	cxt.Response.PackageLength()
	secret := cxt.RadNas.Secret
	replyAuthenticator(cxt.Request.Authenticator, cxt.Response, secret)
	cxt.Listener.WriteToUDP(cxt.Response.ToByte(), cxt.Dst)
	cxt.Next()
}

func AcctRecord(cxt *Context) {
	attr, ok := cxt.Request.RadiusAttrStringKeyMap["User-Name"]
	if !ok {
		cxt.throwPackage = true
		panic("accounting failure, It does not contain User-Name attribute")
	}

	userName := attr.AttrStringValue
	user := RadUser{UserName: userName}
	engine.Get(&user)

	if user.Id == 0 {
		cxt.throwPackage = true
		panic("can not find user " + userName)
	}
	setAcctRecord(userName, cxt)
	cxt.Next()
}

func setAcctRecord(userName string, cxt *Context) {
	attr, _ := cxt.Request.RadiusAttrStringKeyMap["Acct-Status-Type"]
	statusType, _ := strconv.Atoi(attr.AttrStringValue)
	switch statusType {
	case AcctStatusTypeStart:
		acctStartHandler(userName, cxt)
	case AcctStatusTypeStop:
		acctStopHandler(userName, cxt)
	case AcctStatusTypeInterimUpdate:
		acctInterimUpdateHandler(userName, cxt)
	case AcctStatusTypeAccountingOn:
		acctAccountingOn()
	case AcctStatusTypeAccountingOff:
		acctAccountingOff()
	default:
		cxt.throwPackage = true
		panic("radius accounting status type is not supported")
	}
}

func acctStartHandler(userName string, cxt *Context) {
	online := OnlineUser{
		UserName:  userName,
		NasIpAddr: cxt.RadNas.IpAddr,
		StartTime: time.Now(),
	}

	attr, ok := cxt.Request.RadiusAttrStringKeyMap["Acct-Session-Id"]
	if ok {
		online.AcctSessionId = attr.AttrStringValue
	}

	attr, ok = cxt.Request.RadiusAttrStringKeyMap["Framed-IP-Address"]
	if ok {
		online.IpAddr = attr.AttrStringValue
	}

	attr, ok = cxt.Request.RadiusAttrStringKeyMap["NAS-Port-Id"]
	if ok {
		online.NasPortId = attr.AttrStringValue
	}

	online.MacAddr = getMacAddr(cxt)

	_, err := engine.InsertOne(&online)
	if err != nil {
		cxt.throwPackage = true
		panic("insert online user info failure" + err.Error())
	}
}

func acctStopHandler(userName string, cxt *Context) {
	attr, _ := cxt.Request.RadiusAttrStringKeyMap["Acct-Session-Id"]
	online := OnlineUser{AcctSessionId: attr.AttrStringValue}
	engine.Get(&online)

	if online.Id == 0 {
		cxt.throwPackage = true
		panic("in online records can not find this user " + userName)
	}

	// 单位KB
	var totalUpStream, totalDownStream int
	attr, ok := cxt.Request.RadiusAttrStringKeyMap["Acct-Input-Octets"]
	if ok {
		totalUpStream += int(binary.BigEndian.Uint32(attr.AttrValue)) / 1024
	}

	attr, ok = cxt.Request.RadiusAttrStringKeyMap["Acct-Output-Octets"]
	if ok {
		totalDownStream += int(binary.BigEndian.Uint32(attr.AttrValue)) / 1024
	}

	// This attribute indicates how many times the Acct-Input-Octets
	// counter has wrapped around 2^32 over the course of this service being provided
	attr, ok = cxt.Request.RadiusAttrStringKeyMap["Acct-Input-Gigawords"]
	if ok {
		totalUpStream += int(binary.BigEndian.Uint32(attr.AttrValue)) * 4 * 1024 * 1024
	}

	attr, ok = cxt.Request.RadiusAttrStringKeyMap["Acct-Input-Gigawords"]
	if ok {
		totalDownStream += int(binary.BigEndian.Uint32(attr.AttrValue)) * 4 * 1024 * 1024
	}

	// 添加online log
	now := time.Now()
	usedDuration := int(now.Sub(online.StartTime).Seconds())
	onlineLog := UserOnlineLog{
		UserName:        userName,
		StartTime:       online.StartTime,
		StopTime:        &now,
		UsedDuration:    usedDuration,
		TotalUpStream:   totalUpStream,
		TotalDownStream: totalDownStream,
		NasIpAddr:       online.NasIpAddr,
		IpAddr:          online.IpAddr,
		MacAddr:         online.MacAddr,
	}
	engine.InsertOne(&onlineLog)

	// 扣除用户流量，时长
	user := RadUser{UserName: userName}
	engine.Get(&user)
	user.AvailableFlow -= int64(totalDownStream) - int64(totalUpStream)
	user.AvailableTime -= usedDuration
	engine.Cols("available_flow","available_time").Update(&user)

	// 删除online
	delOnline := &OnlineUser{}
	engine.Id(online.Id).Delete(delOnline)
}

func acctInterimUpdateHandler(userName string, cxt *Context) {

}

// It may also be used to mark the start of accounting (for example, upon booting)
// by specifying Accounting-On and to mark the end of accounting
// (for example, just before a scheduled reboot) by specifying Accounting-Off.
func acctAccountingOn() {

}

func acctAccountingOff() {

}
