package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// radius报文结构
type RadiusPackage struct {
	Code byte
	Identifier byte
	// The minimum length is 20 and maximum length is 4096.
	Length uint16
	// radius Authenticator 认证字
	Authenticator [16]byte
	AuthenticatorString string
	// radius attributes
	RadiusAttrs [] RadiusAttr
	//是否是chap请求
	isChap bool
	challenge [16]byte
}

func (r *RadiusPackage) AddRadiusAttr(attr RadiusAttr)  {
	if attr.Length() == 0 {
		return
	}
	r.RadiusAttrs = append(r.RadiusAttrs, attr)
}

// 计算radius package长度
func (r *RadiusPackage) PackageLength() {
	length := PACKAGE_HEADER_LENGTH
	for _, ra :=range r.RadiusAttrs {
		length += int(ra.AttrLength)
	}
	r.Length = uint16(length)
}

// 网络字节序，大端
func (r *RadiusPackage) ToByte() []byte {
	var buf bytes.Buffer
	buf.WriteByte(r.Code)
	buf.WriteByte(r.Identifier)

	var bs = make([]byte,0,  2)
	binary.BigEndian.PutUint16(bs, r.Length)
	buf.Write(bs)

	buf.Write(r.Authenticator[:])
	// radius attr insert into bytes
	return buf.Bytes()
}

// radius attribute
type RadiusAttr struct {
	AttrType byte
	AttrName string
	AttrLength byte
	AttrValue []byte
	AttrStringValue string
	// 26号私有属性专用
	VendorId uint32
	VendorAttrs []VendorAttr
}

func (r RadiusAttr) String1() string {
	return fmt.Sprintf("{%s = %s}", r.AttrName, r.AttrStringValue)
}

func (r *RadiusAttr) setStandardAttrStringVal() {
	attrKey := AttrKey{r.VendorId, int(r.AttrType)}
	attribute, ok := ATTRITUBES[attrKey]
	if ok {
		r.AttrStringValue = getAttrValue(attribute.ValueType, r.AttrValue)
	}
}

func (r *RadiusAttr) Length() byte {
	r.AttrLength = byte(ATTR_HEADER_LENGHT + len(r.AttrValue))

	if r.VendorId != 0 && r.AttrType == VENDOR_SPECIFIC_TYPE {
		var vendorAttrLen byte = 0
		vendorAttrs := r.VendorAttrs
		if len(vendorAttrs) > 0 {
			for _, va := range vendorAttrs {
				vendorAttrLen += va.Length()
			}
		}
		r.AttrLength = ATTR_HEADER_LENGHT + VENDOR_ID_LENGTH + vendorAttrLen
	}

	return r.AttrLength
}


func (r *RadiusAttr) addSpecRadiusAttr(vendorAttr VendorAttr) {
	r.VendorAttrs = append(r.VendorAttrs, vendorAttr)
}

func (r *RadiusAttr) toBytes() []byte {
	bs := make([]byte, 0, r.AttrLength)
	_ = append(bs, r.AttrType, r.AttrLength)
	if r.VendorId != 0 && r.AttrType == VENDOR_SPECIFIC_TYPE {
		var bts []byte
		binary.BigEndian.PutUint32(bts, r.VendorId)
		_ = append(bs, bts...)

		for _, va := range r.VendorAttrs {
			_ = append(bs, va.toBytes()...)
		}

		return bs
	}
	_ = append(bs, r.AttrValue...)
	return bs
}

// radius厂商定义的私有属性
type VendorAttr struct {
	VendorId uint32
	VendorType byte
	VendorTypeName string
	VendorLength byte
	VendorValue []byte
	VendorValueString string
}

func (r VendorAttr) String1() string {
	return fmt.Sprintf("{%s=%s}", r.VendorTypeName, r.VendorValueString)
}

// 获取厂商私有属性长度
func (r *VendorAttr) Length() byte {
	r.VendorLength = byte(len(r.VendorValue) + ATTR_HEADER_LENGHT)
	return r.VendorLength
}

func (r *VendorAttr) toBytes() []byte {
	bs := make([]byte, 0, r.VendorLength)
	_ = append(bs, r.VendorType, r.VendorLength)
	_ = append(bs,  r.VendorValue...)
	return bs
}

func (r *VendorAttr) setVendorAttrStringValue() {
	attribute, ok := ATTRITUBES[AttrKey{r.VendorId, int(r.VendorType)}]
	if ok {
		r.VendorValueString = getAttrValue(attribute.ValueType, r.VendorValue)
	}
}