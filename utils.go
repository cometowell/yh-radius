package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)


func LeftPadChar(source string, padChar byte, size int) string {
	sourceLength := len(source)
	if sourceLength >= size {
		return source
	}
	return strings.Repeat(string(padChar), size - sourceLength) + source
}

func RightPadChar(source string, padChar byte, size int) string {
	sourceLength := len(source)
	if sourceLength >= size {
		return source
	}
	return source + strings.Repeat(string(padChar), size - sourceLength)
}

func IpAddrToBytes(ipAddr string) (ipArr []byte, err error) {
	ipArr = make([]byte, 4)
	items := strings.Split(ipAddr, ".")
	if len(items) != 4 {
		return nil, errors.New("ip地址格式错误")
	}

	for index, item := range items {
		val, e := strconv.Atoi(item)
		if e != nil {
			return nil, e
		}
		ipArr[index] = byte(val)
	}

	return ipArr, nil
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

func FillBytesByString(size int, value string) []byte {
	if len(value) >= size {
		return []byte(value)
	}
	ret := make([]byte, size)
	copy(ret, []byte(value))
	return ret
}

func getIntegerBytes(val uint32) []byte {
	container := make([]byte, 4)
	binary.BigEndian.PutUint32(container, val)
	return container
}

func setVendorStringValue(vendorId uint32, vendorAttr *VendorAttr) {
	attr, ok := ATTRITUBES[AttrKey{vendorId, int(vendorAttr.VendorType)}]
	if ok {
		vendorAttr.VendorId = vendorId
		vendorAttr.VendorTypeName = attr.Name
		vendorAttr.VendorValueString = getAttrValue(attr.ValueType, vendorAttr.VendorValue)
	}
}