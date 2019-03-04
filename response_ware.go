package main

import "encoding/binary"

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
	product := RadProduct{}
	engine.Id(user.ProductId).Get(&product)
	// Huawei-Input-Burst-Size, Huawei-Input-Average-Rate
	// Huawei-Output-Burst-Size, Huawei-Output-Average-Rate 单位 bit/s
	upStreamLimit := product.UpStreamLimit * 1024 * 8
	downStreamLimit := product.DownStreamLimit * 1024 * 8
	specAttr := &RadiusAttr{
		AttrType: VendorSpecificType,
		VendorId: Huawei,
		AttrName: "Vendor-Specific",
	}

	container := make([]byte, 4)
	// 上行最大速率
	binary.BigEndian.PutUint32(container, upStreamLimit)
	inputBurstVendorAttr := VendorAttr{
		VendorType:   1,
		VendorLength: 6,
		VendorValue:  container,
	}
	setVendorStringValue(Huawei, &inputBurstVendorAttr)
	specAttr.addSpecRadiusAttr(inputBurstVendorAttr)


	// 上行平均速率
	binary.BigEndian.PutUint32(container, upStreamLimit)
	inputAverageVendorAttr := VendorAttr{
		VendorType:   2,
		VendorLength: 6,
		VendorValue:  container,
	}
	setVendorStringValue(Huawei, &inputAverageVendorAttr)
	specAttr.addSpecRadiusAttr(inputAverageVendorAttr)

	// 下行最大速率
	binary.BigEndian.PutUint32(container, downStreamLimit)
	outputBurstVendorAttr := VendorAttr{
		VendorType:   3,
		VendorLength: 6,
		VendorValue:  container,
	}
	setVendorStringValue(Huawei, &outputBurstVendorAttr)
	specAttr.addSpecRadiusAttr(outputBurstVendorAttr)

	// 下行平均速率
	binary.BigEndian.PutUint32(container, downStreamLimit)
	outputAverageVendorAttr := VendorAttr{
		VendorType:   4,
		VendorLength: 6,
		VendorValue:  container,
	}
	setVendorStringValue(Huawei, &outputAverageVendorAttr)
	specAttr.addSpecRadiusAttr(outputAverageVendorAttr)

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

func setVendorStringValue(vendorId uint32, vendorAttr *VendorAttr) {
	attr, ok := ATTRITUBES[AttrKey{vendorId, int(vendorAttr.VendorType)}]
	if ok {
		vendorAttr.VendorId = vendorId
		vendorAttr.VendorTypeName = attr.Name
		vendorAttr.VendorValueString = getAttrValue(attr.ValueType, vendorAttr.VendorValue)
	}
}

// 思科
func CiscoResponse(cxt *Context) {

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
