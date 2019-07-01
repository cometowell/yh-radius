package radius

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"go-rad/common"
	"go-rad/database"
	"go-rad/model"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func SetVendorStringValue(vendorId uint32, vendorAttr *VendorAttr) {
	attr, ok := ATTRITUBES[AttrKey{vendorId, int(vendorAttr.VendorType)}]
	if ok {
		vendorAttr.VendorId = vendorId
		vendorAttr.VendorTypeName = attr.Name
		vendorAttr.VendorValueString = getAttrValue(attr.ValueType, vendorAttr.VendorValue)
	}
}

func buildUrlParamsFromMap(params map[string]interface{}) string {
	var result string = "?"
	for key, val := range params {
		result += fmt.Sprintf("%s=%s&", key, val)
	}

	if strings.HasSuffix(result, "&") {
		result = result[:len(result)-1]
	}

	return result
}

func buildUrlParams(params ...interface{}) string {

	if len(params)%2 != 0 {
		return ""
	}

	var result string = "?"
	for i := 0; i < len(params); i += 2 {
		result += fmt.Sprintf("%s=%s&", params[i], params[i+1])
	}

	if strings.HasSuffix(result, "&") {
		result = result[:len(result)-1]
	}

	return result
}

func IsExpire(user *model.RadUser, product *model.RadProduct) bool {
	duration := time.Time(user.ExpireTime).Sub(time.Now())
	if duration <= 0 {
		return true
	} else if product.Type == common.TimeProduct {
		return user.AvailableTime <= 0
	} else if product.Type == common.FlowProduct {
		return user.AvailableFlow <= 0
	}

	return false
}

// offline user
// request Authenticator define
// The NAS and RADIUS accounting server share a secret.  The Request
// Authenticator field in Accounting-Request packets contains a one-
// way MD5 hash calculated over a stream of octets consisting of the
// Code + Identifier + Length + 16 zero octets + request attributes +
// shared secret (where + indicates concatenation).  The 16 octet MD5
// hash value is stored in the Authenticator field of the
// Accounting-Request packet.

// now, only send the dm request don't handler the response from nas
// when user offline, nas will send stop accounting package to radius server
// then radius handler the stop accounting package
// maybe After sometime, add return response package processing
func OfflineUser(online model.RadOnlineUser) error {

	var nas model.RadNas
	ok, _ := database.DataBaseEngine.Where("ip_addr = ?", online.NasIpAddr).Get(&nas)
	if !ok {
		return errors.New("nas can not be found")
	}

	rp := RadiusPackage{}
	rp.Authenticator = [16]byte{}
	rp.Code = DisconnectRequest
	rp.Identifier = byte(rand.Intn(256))

	attrs := make([]*RadiusAttr, 0, 3)
	acctSessionIdAttr := RadiusAttr{
		AttrType:  44,
		AttrValue: []byte(online.AcctSessionId),
	}
	acctSessionIdAttr.Length()
	attrs = append(attrs, &acctSessionIdAttr)
	ipArrBytes, _ := common.IpAddrToBytes(online.NasIpAddr)
	nasIpAddrAttr := RadiusAttr{
		AttrType:   4,
		AttrLength: 6,
		AttrValue:  ipArrBytes,
	}
	attrs = append(attrs, &nasIpAddrAttr)

	if online.UserName != "" {
		usernameAttr := RadiusAttr{
			AttrType:  1,
			AttrValue: []byte(online.UserName),
		}
		usernameAttr.Length()
		attrs = append(attrs, &usernameAttr)
	}

	rp.RadiusAttrs = attrs
	replyAuthenticator(rp.Authenticator, &rp, nas.Secret)
	rp.PackageLength()

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", online.NasIpAddr, nas.AuthorizePort))
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(rp.ToByte())
	if err != nil {
		return err
	}

	return nil
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

// ResponseAuth = MD5(Code+ID+Length+RequestAuth+Attributes+Secret)
func replyAuthenticator(authAuthenticator [16]byte, reply *RadiusPackage, secret string) {
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
				ciscoMacAddr := strings.Replace(matchs[1], ".", "", -1)
				return strings.ToUpper(fmt.Sprintf("%s:%s:%s:%s:%s:%s",
					ciscoMacAddr[0:2], ciscoMacAddr[2:4], ciscoMacAddr[4:6],
					ciscoMacAddr[6:8], ciscoMacAddr[8:10], ciscoMacAddr[10:12]),
				)
			}
		}
	}
	return ""
}
