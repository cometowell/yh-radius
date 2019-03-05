package main

import (
	"encoding/binary"
	"fmt"
)

// radius认证响应处理handlers
// 根据不同的厂商需要下发不同的属性：限速，限定域，会话时长及其他
// key = vendorId, value = middleware
var AuthResponseWares = map[int]func(cxt *Context){
	Standard: StandardResponse,
	Cisco:    CiscoResponse,
	Huawei:   HuaweiResponse,
	Zte:      zteResponse,
	MikroTik: MikroTikResponse,
}

// 华为
func HuaweiResponse(cxt *Context) {
	user := cxt.User
	product := user.product
	// Huawei-Input-Burst-Size, Huawei-Input-Average-Rate
	// Huawei-Output-Burst-Size, Huawei-Output-Average-Rate 单位 bit/s
	upStreamLimit := product.UpStreamLimit * 1024 * 8
	downStreamLimit := product.DownStreamLimit * 1024 * 8
	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: Huawei,
		AttrName: "Vendor-Specific",
	}

	inputAvgRateVendorAttr := VendorAttr{
		VendorType:   2,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(upStreamLimit),
	}
	setVendorStringValue(Huawei, &inputAvgRateVendorAttr)
	specAttr.addSpecRadiusAttr(inputAvgRateVendorAttr)


	// 上行平均速率
	inputPeakRateVendorAttr := VendorAttr{
		VendorType:   3,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(upStreamLimit),
	}
	setVendorStringValue(Huawei, &inputPeakRateVendorAttr)
	specAttr.addSpecRadiusAttr(inputPeakRateVendorAttr)

	// 下行最大速率
	outputAvgRateVendorAttr := VendorAttr{
		VendorType:   5,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(downStreamLimit),
	}
	setVendorStringValue(Huawei, &outputAvgRateVendorAttr)
	specAttr.addSpecRadiusAttr(outputAvgRateVendorAttr)

	// 下行平均速率
	outputPeakRateVendorAttr := VendorAttr{
		VendorType:   6,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(downStreamLimit),
	}
	setVendorStringValue(Huawei, &outputPeakRateVendorAttr)
	specAttr.addSpecRadiusAttr(outputPeakRateVendorAttr)

	// context Huawei-Domain-Name 64字节
	if product.DomainName != "" {
		domainNameVendorAttr := VendorAttr{
			VendorType:   138,
			VendorLength: 66,
			VendorValue:  FillBytesByString(64, product.DomainName),
		}
		setVendorStringValue(Huawei, &domainNameVendorAttr)
		specAttr.addSpecRadiusAttr(domainNameVendorAttr)
	}

	specAttr.Length()

	cxt.Response.AddRadiusAttr(*specAttr)
}


// 思科
func CiscoResponse(cxt *Context) {
	product := cxt.User.product
	upStreamLimit := product.UpStreamLimit * 1024 * 8
	downStreamLimit := product.DownStreamLimit * 1024 * 8

	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: Cisco,
		AttrName: "Vendor-Specific",
	}

	qosIn := VendorAttr{
		VendorType: 1,
		VendorValue: []byte(fmt.Sprintf("sub-qos-policy-in=%d", upStreamLimit)),
	}
	qosIn.Length()
	setVendorStringValue(Cisco, &qosIn)
	specAttr.addSpecRadiusAttr(qosIn)

	qosOut := VendorAttr{
		VendorType: 1,
		VendorValue: []byte(fmt.Sprintf("sub-qos-policy-out=%d", downStreamLimit)),
	}
	qosOut.Length()
	setVendorStringValue(Cisco, &qosOut)
	specAttr.addSpecRadiusAttr(qosOut)

	if product.DomainName != "" {
		domainNameVendorAttr := VendorAttr{
			VendorType:   1,
			VendorValue:  []byte(fmt.Sprintf("addr-pool=%s", product.DomainName)),
		}
		domainNameVendorAttr.Length()
		setVendorStringValue(Cisco, &domainNameVendorAttr)
		specAttr.addSpecRadiusAttr(domainNameVendorAttr)
	}

	specAttr.Length()
	cxt.Response.AddRadiusAttr(*specAttr)
}

// RFC标准
func StandardResponse(cxt *Context) {

}

// RouterOS
func MikroTikResponse(cxt *Context) {

}

// 中兴
func zteResponse(cxt *Context) {

}

// 设置认证响应属性
func AuthSpecAndCommonAttrSetter(cxt *Context) {
	vendorId := cxt.RadNas.VendorId
	vendorRespFunc, ok := AuthResponseWares[vendorId]
	if ok {
		vendorRespFunc(cxt)
	}

	// session timeout
	sessionTimeoutAttr := RadiusAttr{
		AttrType: 27,
		// 默认会话时长一星期
		AttrValue: getIntegerBytes(uint32(cxt.User.sessionTimeout)),
	}
	sessionTimeoutAttr.Length()
	attr, _ := ATTRITUBES[AttrKey{Standard, int(sessionTimeoutAttr.AttrType)}]
	sessionTimeoutAttr.AttrName = attr.Name

	//TODO 设置用户会话时长

	sessionTimeoutAttr.setStandardAttrStringVal()
	cxt.Response.AddRadiusAttr(sessionTimeoutAttr)
	cxt.Next()
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