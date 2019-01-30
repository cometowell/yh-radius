package main

import (
	"bytes"
	"encoding/binary"
)

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
}

// 注意RadiusAttr length == 0不加入到属性列表
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

// 网络字节序，大端方式
func (r *RadiusPackage) ToByte() []byte {
	var buf bytes.Buffer
	buf.WriteByte(r.Code)
	buf.WriteByte(r.Identifier)

	var bs = make([]byte, 2)
	binary.BigEndian.PutUint16(bs, r.Length)
	buf.Write(bs)

	buf.Write(r.Authenticator[:])
	// radius attr insert into bytes
	return buf.Bytes()
}

// radius attribute
type RadiusAttr struct {
	AttrType byte
	AttrLength byte
	AttrValue []byte

	// 26号私有属性专用
	VendorId uint32
	// RFC定义中为string类型
	VendorAttrs []VendorAttr
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
	bs := make([]byte, r.AttrLength)
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
	VendorType byte
	VendorLength byte
	VendorValue []byte
}

// 获取厂商私有属性长度
func (r *VendorAttr) Length() byte {
	r.VendorLength = byte(len(r.VendorValue) + ATTR_HEADER_LENGHT)
	return r.VendorLength
}

func (r *VendorAttr) toBytes() []byte {
	bs := make([]byte, r.VendorLength)
	_ = append(bs, r.VendorType, r.VendorLength)
	_ = append(bs,  r.VendorValue...)
	return bs
}