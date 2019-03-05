package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	AccessAcceptReplyMsg = "authenticate success"
	ShouldBindMacAddrFlag = 1 // 用户绑定MAC地址标识
	ShouldBindVlanFlag 		  // 用户绑定虚拟专网标识
)

func NasValidation(cxt *Context) {
	nasIp := cxt.Dst.IP.String()
	logger.Infoln("UDP报文消息来源：", nasIp)
	logger.Infof("%v\n", cxt.Request)

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
		panic("user's account number or password is incorrect")
	}

	userName := attr.AttrStringValue
	user := RadUser{UserName: userName}
	engine.Get(&user)

	if user.Id == 0 {
		panic("user's account number or password is incorrect")
	}

	// 验证密码
	password := decrypt(user.Password)
	if cxt.Request.isChap {
		if !chap(password, &cxt.Request) {
			panic("user's account number or password is incorrect")
		}
	} else {
		if !pap(cxt.RadNas.Secret, password, cxt.Request) {
			panic("user's account number or password is incorrect")
		}
	}

	onlineUser := OnlineUser{UserName: userName}
	onlineCount, _ := engine.Count(onlineUser)
	userConcurrent := user.ConcurrentCount
	if userConcurrent != 0 && userConcurrent <= uint(onlineCount) {
		panic(fmt.Sprintf("the maximum number of concurrency has been reached: %d", onlineCount))
	}

	product := RadProduct{}
	engine.Id(user.ProductId).Get(&product)

	if product.Id == 0 {
		panic("user did not purchase the product")
	}

	user.product = product
	productType := product.Type

	sessionTimeout := int(config["radius.session.timeout"].(float64))
	user.sessionTimeout = sessionTimeout
	if productType == ClacPriceByMonth { // 按月计费套餐
		availableSeconds := int(user.ExpireTime.Sub(time.Now()).Seconds())
		if availableSeconds <= 0 {
			panic("user's service is expire")
		}
		if sessionTimeout > availableSeconds {
			user.sessionTimeout = availableSeconds
		}

	} else if productType == UseTimesProductType { // 时长套餐
		if user.AvailableTime <= 0 {
			panic("user's service time already used up")
		}
		if sessionTimeout > user.AvailableTime {
			user.sessionTimeout = user.AvailableTime
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
			engine.Id(user.Id).Cols("mac_addr").Update(user)
		}

		if macAddr != user.MacAddr {
			logger.Panicf("用户MAC地址: %s != %s", user.MacAddr, macAddr)
		}
	}
	cxt.Next()
}

// 获取MAC地址，使用：分隔方式 AA:BB:CC:DD:EE:FF
// 有些厂商的MAC地址需要从(type=26)私有属性中获取
func getMacAddr(cxt *Context) string {
	vendorId := cxt.RadNas.VendorId
	if vendorId == Standard {
		attr, ok := cxt.Request.RadiusAttrStringKeyMap["Calling-Station-Id"]
		if ok {
			return strings.ToUpper(attr.AttrStringValue)
		}
	}
	return getVendorMacAddr(vendorId, cxt)
}

func getVendorMacAddr(vendorId int, cxt *Context) string {
	if vendorId == Huawei {
		attr, ok := cxt.Request.RadiusAttrStringKeyMap["Calling-Station-Id"]
		if ok {
			return strings.ToUpper(attr.AttrStringValue)
		}
	} else if vendorId == Cisco {
		avPairParttern, _ := regexp.Compile(`client-mac-address=(\w{4}\.\w{4}\.\w{4})`)
		attr, ok := cxt.Request.RadiusAttrStringKeyMap["Vendor-Specific"]
		if ok {
			ciscoAVPair, ok := attr.VendorAttrStringKeyMap["Cisco-AVPair"]
			if ok {
				avPairVal := ciscoAVPair.VendorValueString
				matchs := avPairParttern.FindStringSubmatch(avPairVal)
				ciscoMacAddr := strings.Replace(matchs[1],".","", -1)
				return strings.ToUpper(fmt.Sprintf("%s:%s:%s:%s:%s:%s",
					ciscoMacAddr[0:2],ciscoMacAddr[2:4],ciscoMacAddr[4:6],
					ciscoMacAddr[6:8],ciscoMacAddr[8:10],ciscoMacAddr[10:12]),
				)
			}
		}
	}
	return ""
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
				logger.Error(msg)
				panic(msg)
			}

			if shouldUpdate {
				engine.Id(user.Id).Cols("vlan_id", "vlan_id2").Update(user)
			}
		}
	}
	cxt.Next()
}

// 不同厂商不同的解析方式，这里是通用的方式
func getVlanIds(vendorId int, nasPortId string) (int, int) {
	var ptn *regexp.Regexp
	var retMatch []string
	if vendorId == Cisco {
		// eth phy_slot/phy_subslot/phy_port:XPI.XCI
		ptn, _ = regexp.Compile(`phy_port:(\d).(\d)`)
		retMatch = ptn.FindStringSubmatch(nasPortId)
	} else {
		ptn, _ = regexp.Compile(`vlanid=(\d);vlanid2=(\d)`)
		retMatch = ptn.FindStringSubmatch(nasPortId)
	}

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
					var errMsg string
					if entry, ok := err.(*logrus.Entry); ok {
						errMsg = entry.Message
					} else if msg, ok := err.(string); ok {
						errMsg = msg
					}  else {
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
	reply.AuthenticatorString = hex.EncodeToString(sum)
}
