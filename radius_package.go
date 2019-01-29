package main

// radius报文结构
type RadiusPackage struct {
	Code byte
	Identifier byte
	// The minimum length is 20 and maximum length is 4096.
	Length uint16
	// radius Authenticator 认证字
	Authenticator [16]byte
	// radius attributes slice
	RadiusAttrs []RadiusAttr
	// 26 vendor specific 单独处理
	RadiusVendorSpecificAttr RadiusVendorSpecificAttr
}

func (r *RadiusPackage) AddRadiusAttr(attr RadiusAttr)  {
	r.RadiusAttrs = append(r.RadiusAttrs, attr)
}

func (r *RadiusPackage) SetRadiusVendorSpecificAttr(vendorAttrs []VendorAttr)  {
	r.RadiusVendorSpecificAttr.AttrType = VENDOR_SPECIFIC_TYPE
	r.RadiusVendorSpecificAttr.VendorAttrs = append(r.RadiusVendorSpecificAttr.VendorAttrs, vendorAttrs...)
	// 计算长度length
	r.RadiusVendorSpecificAttr.AttrLength = 0
}

// radius attribute
type RadiusAttr struct {
	AttrType byte
	AttrLength byte
	AttrValue []byte
}

// radius 26 Vendor-Specific
type RadiusVendorSpecificAttr struct {
	AttrType byte
	AttrLength byte
	VendorId uint32
	// RFC定义中为string类型
	VendorAttrs []VendorAttr
}

func (r *RadiusVendorSpecificAttr) addSpecRadiusAttr(vendorAttr VendorAttr) {
	r.VendorAttrs = append(r.VendorAttrs, vendorAttr)
}

// radius厂商定义的私有属性
type VendorAttr struct {
	VendorType byte
	VendorLength byte
	VendorValue []byte
}