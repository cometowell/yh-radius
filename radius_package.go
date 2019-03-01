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
	RadiusAttrMap map[AttrKey]RadiusAttr `json:"-"`
	RadiusAttrStringKeyMap map[string]RadiusAttr `json:"-"`
	//是否是chap请求
	isChap bool
	challenge []byte

}

func (r RadiusPackage) String() string {
	return fmt.Sprintf(`RadiusPackage:{
        Code=%d
		Identifier=%d
		Length=%d
		Authenticator=%s
		isChap=%v
		RadiusAttrs: %v
	}`, r.Code, r.Identifier, r.Length, r.AuthenticatorString, r.isChap, r.RadiusAttrs)
}

func (r *RadiusPackage) AddRadiusAttr(attr RadiusAttr)  {
	if attr.Length() == 0 {
		return
	}
	r.RadiusAttrs = append(r.RadiusAttrs, attr)
}

// 计算radius package长度
func (r *RadiusPackage) PackageLength() {
	length := PackageHeaderLength
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

	var attrBuf bytes.Buffer
	for _, attr := range r.RadiusAttrs {
		attrBuf.Write(attr.toBytes())
	}

	r.PackageLength()

	var bs = make([]byte,2)
	binary.BigEndian.PutUint16(bs, r.Length)
	buf.Write(bs)

	buf.Write(r.Authenticator[:])
	buf.Write(attrBuf.Bytes())

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
	VendorAttrMap map[AttrKey]VendorAttr
	VendorAttrStringKeyMap map[string]VendorAttr
}

func (r RadiusAttr) String() string {
	if r.AttrType != VendorSpecificType {
		return fmt.Sprintf("\n \t\t\t %s=%s", r.AttrName, r.AttrStringValue)
	} else {
		return fmt.Sprintf("\n \t\t\t VendorAttrs: %s", r.VendorAttrs)
	}
}

func (r *RadiusAttr) setStandardAttrStringVal() {
	attrKey := AttrKey{r.VendorId, int(r.AttrType)}
	attribute, ok := ATTRITUBES[attrKey]
	if ok {
		r.AttrStringValue = getAttrValue(attribute.ValueType, r.AttrValue)
	}
}

func (r *RadiusAttr) Length() byte {
	r.AttrLength = byte(AttrHeaderLength + len(r.AttrValue))

	if r.VendorId != 0 && r.AttrType == VendorSpecificType {
		var vendorAttrLen byte = 0
		vendorAttrs := r.VendorAttrs
		if len(vendorAttrs) > 0 {
			for _, va := range vendorAttrs {
				vendorAttrLen += va.Length()
			}
		}
		r.AttrLength = AttrHeaderLength + VendorIdLength + vendorAttrLen
	}

	return r.AttrLength
}


func (r *RadiusAttr) addSpecRadiusAttr(vendorAttr VendorAttr) {
	r.VendorAttrs = append(r.VendorAttrs, vendorAttr)
}

func (r *RadiusAttr) toBytes() []byte {
	bs := make([]byte, 0, r.AttrLength)
	bs = append(bs, r.AttrType, r.AttrLength)
	if r.VendorId != 0 && r.AttrType == VendorSpecificType {
		var bts []byte
		binary.BigEndian.PutUint32(bts, r.VendorId)
		bs = append(bs, bts...)

		for _, va := range r.VendorAttrs {
			bs = append(bs, va.toBytes()...)
		}

		return bs
	}
	bs = append(bs, r.AttrValue...)
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

func (r VendorAttr) String() string {
	return fmt.Sprintf("\n \t\t\t   %s=%s", r.VendorTypeName, r.VendorValueString)
}

// 获取厂商私有属性长度
func (r *VendorAttr) Length() byte {
	r.VendorLength = byte(len(r.VendorValue) + AttrHeaderLength)
	return r.VendorLength
}

func (r *VendorAttr) toBytes() []byte {
	bs := make([]byte, 0, r.VendorLength)
	bs = append(bs, r.VendorType, r.VendorLength)
	bs = append(bs,  r.VendorValue...)
	return bs
}

func (r *VendorAttr) setVendorAttrStringValue() {
	attribute, ok := ATTRITUBES[AttrKey{r.VendorId, int(r.VendorType)}]
	if ok {
		r.VendorValueString = getAttrValue(attribute.ValueType, r.VendorValue)
	}
}