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



// 获得26号私有属性的长度
func (r *RadiusVendorSpecificAttr) getLength() int {
	vendorAttrLen := 0
	vendorAttrs := r.VendorAttrs
	if len(vendorAttrs) > 0 {
		for _, va := range vendorAttrs {
			vendorAttrLen += va.getLength()
		}
	}
	return ATTR_TYPE_FIELD_LENGHT + ATTR_LENGTH_FIELD_LENGHT + VENDOR_ID_LENGTH + vendorAttrLen
}

// radius厂商定义的私有属性
type VendorAttr struct {
	VendorType byte
	VendorLength byte
	VendorValue []byte
}

// 获取厂商私有属性长度
func (r *VendorAttr) getLength() int {
	return len(r.VendorValue) + ATTR_TYPE_FIELD_LENGHT + ATTR_LENGTH_FIELD_LENGHT
}