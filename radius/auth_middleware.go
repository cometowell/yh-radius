package radius

import (
	"fmt"
	"go-rad/common"
	"go-rad/logger"
	"go-rad/model"
	"time"
)

const (
	AccessAcceptReplyMsg  = "authenticate success"
	ShouldBindMacAddrFlag = 1 // 用户绑定MAC地址标识
	ShouldBindVlanFlag        // 用户绑定虚拟专网标识
)

// 验证用户名，密码
func UserVerify(cxt *Context) {
	attr, ok := cxt.Request.RadiusAttrStringKeyMap["User-Name"]
	if !ok {
		panic("user's account number or password is incorrect")
	}

	userName := attr.AttrStringValue
	user := model.RadUser{UserName: userName}
	cxt.Session.Get(&user)

	if user.Id == 0 {
		panic("user's account number or password is incorrect")
	}

	if user.Status > UserAvailableStatus {
		panic("status of user's account number status abnormal")
	}

	// 验证密码
	password := common.Decrypt(user.Password)
	if cxt.Request.isChap {
		if !chap(password, &cxt.Request) {
			panic("user's account number or password is incorrect")
		}
	} else {
		if !pap(cxt.RadNas.Secret, password, cxt.Request) {
			panic("user's account number or password is incorrect")
		}
	}

	onlineUser := model.RadOnlineUser{UserName: userName}
	onlineCount, _ := cxt.Session.Count(onlineUser)
	userConcurrent := user.ConcurrentCount
	if userConcurrent != 0 && userConcurrent <= int(onlineCount) {
		panic(fmt.Sprintf("the maximum number of concurrency has been reached: %d", onlineCount))
	}

	product := model.RadProduct{}
	cxt.Session.ID(user.ProductId).Get(&product)

	if product.Id == 0 {
		panic("user did not purchase the product")
	}

	user.Product = product
	productType := product.Type

	sessionTimeout := int(common.Config["radius.session.timeout"].(float64))
	user.SessionTimeout = sessionTimeout
	if productType == ClacPriceByMonth { // 按月计费套餐
		availableSeconds := int(time.Time(user.ExpireTime).Sub(time.Now()).Seconds())
		if availableSeconds <= 0 {
			panic("user's service is expire")
		}
		if sessionTimeout > availableSeconds {
			user.SessionTimeout = availableSeconds
		}

	} else if productType == UseTimesProductType { // 时长套餐
		if user.AvailableTime <= 0 {
			panic("user's service time already used up")
		}
		if int64(sessionTimeout) > user.AvailableTime {
			user.SessionTimeout = int(user.AvailableTime)
		}
	} else { // 流量套餐
		if user.AvailableFlow <= 0 {
			panic("user's service flow already used up")
		}
	}

	cxt.User = &user
	cxt.Next()
}

// 验证MAC地址绑定
func MacAddrVerify(cxt *Context) {
	user := cxt.User
	if user.ShouldBindMacAddr == ShouldBindMacAddrFlag {
		macAddr := getMacAddr(cxt)
		if user.MacAddr == "" {
			user.MacAddr = macAddr
			cxt.Session.ID(user.Id).Cols("mac_addr").Update(user)
		}

		if macAddr != user.MacAddr {
			logger.Logger.Panicf("用户MAC地址: %s != %s", user.MacAddr, macAddr)
		}
	}
	cxt.Next()
}

// 验证VLAN
func VlanVerify(cxt *Context) {
	user := cxt.User
	if user.ShouldBindVlan == ShouldBindVlanFlag {
		attr, ok := cxt.Request.RadiusAttrStringKeyMap["NAS-Port-Id"]

		if ok {
			vlanId, vlanId2 := getVlanIds(cxt.RadNas.VendorId, attr.AttrStringValue)

			var shouldUpdate bool
			if user.VlanId == 0 && user.VlanId2 == 0 {
				user.VlanId = vlanId
				user.VlanId2 = vlanId2
				shouldUpdate = true
			}

			if vlanId != user.VlanId || vlanId2 != user.VlanId2 {
				msg := fmt.Sprintf("VLAN验证失败用户绑定Vlan信息(VlanId:%d, VlanId2:%d) != (VlanId:%d, VlanId2:%d)", user.VlanId, user.VlanId2, vlanId, vlanId2)
				logger.Logger.Error(msg)
				panic(msg)
			}

			if shouldUpdate {
				cxt.Session.Id(user.Id).Cols("vlan_id", "vlan_id2").Update(user)
			}
		}
	}
	cxt.Next()
}

func AuthAcceptReply(cxt *Context) {
	cxt.Next()
	authReply(cxt, AccessAcceptCode, AccessAcceptReplyMsg)
}

func authReply(cxt *Context, replyCode byte, msg string) {
	cxt.Response.Code = replyCode
	replyMessage := RadiusAttr{
		AttrType:  18,
		AttrValue: []byte(msg),
		VendorId:  Standard,
	}
	replyMessage.Length()
	replyMessage.setStandardAttrStringVal()
	attr, _ := ATTRITUBES[AttrKey{Standard, int(replyMessage.AttrType)}]
	replyMessage.AttrName = attr.Name
	cxt.Response.AddRadiusAttr(replyMessage)
	cxt.Response.PackageLength()
	secret := cxt.RadNas.Secret
	replyAuthenticator(cxt.Request.Authenticator, cxt.Response, secret)
	cxt.Listener.WriteToUDP(cxt.Response.ToByte(), cxt.Dst)
}
