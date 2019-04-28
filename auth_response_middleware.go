package main

import (
	"fmt"
)

// radius认证响应处理handlers
// 根据不同的厂商需要下发不同的属性：限速，限定域，会话时长及其他
// key = vendorId, value = middleware
var AuthResponseWares = map[int]func(cxt *Context){
	Standard: StandardResponse,
	Cisco:    CiscoResponse,
	Huawei:   HuaweiResponse,
	Zte:      ZteResponse,
	MikroTik: MikroTikResponse,
}

// 华为
func HuaweiResponse(cxt *Context) {
	user := cxt.User
	product := user.Product
	// Huawei-Input-Burst-Size, Huawei-Input-Average-Rate
	// Huawei-Output-Burst-Size, Huawei-Output-Average-Rate 单位 bit/s
	upStreamLimit := product.UpStreamLimit * 1024 * 1024
	downStreamLimit := product.DownStreamLimit * 1024 * 1024
	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: Huawei,
		AttrName: "Vendor-Specific",
	}

	inputAvgRateVendorAttr := VendorAttr{
		VendorType:   2,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(uint32(upStreamLimit)),
	}
	setVendorStringValue(Huawei, &inputAvgRateVendorAttr)
	specAttr.addSpecRadiusAttr(inputAvgRateVendorAttr)

	// 上行平均速率
	inputPeakRateVendorAttr := VendorAttr{
		VendorType:   3,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(uint32(upStreamLimit)),
	}
	setVendorStringValue(Huawei, &inputPeakRateVendorAttr)
	specAttr.addSpecRadiusAttr(inputPeakRateVendorAttr)

	// 下行最大速率
	outputAvgRateVendorAttr := VendorAttr{
		VendorType:   5,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(uint32(downStreamLimit)),
	}
	setVendorStringValue(Huawei, &outputAvgRateVendorAttr)
	specAttr.addSpecRadiusAttr(outputAvgRateVendorAttr)

	// 下行平均速率
	outputPeakRateVendorAttr := VendorAttr{
		VendorType:   6,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(uint32(downStreamLimit)),
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
	product := cxt.User.Product
	upStreamLimit := product.UpStreamLimit * 1024 * 1024
	downStreamLimit := product.DownStreamLimit * 1024 * 1024

	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: Cisco,
		AttrName: "Vendor-Specific",
	}

	qosIn := VendorAttr{
		VendorType:  1,
		VendorValue: []byte(fmt.Sprintf("sub-qos-policy-in=%d", upStreamLimit)),
	}
	qosIn.Length()
	setVendorStringValue(Cisco, &qosIn)
	specAttr.addSpecRadiusAttr(qosIn)

	qosOut := VendorAttr{
		VendorType:  1,
		VendorValue: []byte(fmt.Sprintf("sub-qos-policy-out=%d", downStreamLimit)),
	}
	qosOut.Length()
	setVendorStringValue(Cisco, &qosOut)
	specAttr.addSpecRadiusAttr(qosOut)

	if product.DomainName != "" {
		domainNameVendorAttr := VendorAttr{
			VendorType:  1,
			VendorValue: []byte(fmt.Sprintf("addr-pool=%s", product.DomainName)),
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
	// Mikrotik-Rate-Limit	8
	product := cxt.User.Product
	upStreamLimit := product.UpStreamLimit * 1024
	downStreamLimit := product.DownStreamLimit * 1024

	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: MikroTik,
		AttrName: "Vendor-Specific",
	}

	rateLimitAttr := VendorAttr{
		VendorType:  8,
		VendorValue: []byte(fmt.Sprintf("%dk/%dk", upStreamLimit, downStreamLimit)),
	}
	rateLimitAttr.Length()
	setVendorStringValue(MikroTik, &rateLimitAttr)
	specAttr.addSpecRadiusAttr(rateLimitAttr)

	specAttr.Length()
	cxt.Response.AddRadiusAttr(*specAttr)
}

// 中兴, 限速单位 kbit/s
func ZteResponse(cxt *Context) {
	// ZTE-Rate-Ctrl-SCR-Down	83
	// ZTE-Rate-Ctrl-SCR-Up		89
	product := cxt.User.Product
	upStreamLimit := product.UpStreamLimit
	downStreamLimit := product.DownStreamLimit

	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: Zte,
		AttrName: "Vendor-Specific",
	}

	upRateAttr := VendorAttr{
		VendorType:   89,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(uint32(upStreamLimit)),
	}
	setVendorStringValue(Zte, &upRateAttr)
	specAttr.addSpecRadiusAttr(upRateAttr)

	downRateAttr := VendorAttr{
		VendorType:   83,
		VendorLength: 6,
		VendorValue:  getIntegerBytes(uint32(downStreamLimit)),
	}
	setVendorStringValue(Zte, &downRateAttr)
	specAttr.addSpecRadiusAttr(downRateAttr)

	specAttr.Length()
	cxt.Response.AddRadiusAttr(*specAttr)
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
		VendorId: Standard,
		AttrType: 27,
		// 默认会话时长一星期
		AttrValue: getIntegerBytes(uint32(cxt.User.sessionTimeout)),
	}
	sessionTimeoutAttr.Length()
	attr, _ := ATTRITUBES[AttrKey{Standard, int(sessionTimeoutAttr.AttrType)}]
	sessionTimeoutAttr.AttrName = attr.Name
	sessionTimeoutAttr.setStandardAttrStringVal()
	cxt.Response.AddRadiusAttr(sessionTimeoutAttr)
	cxt.Next()
}
