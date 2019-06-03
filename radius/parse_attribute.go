package radius

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const ATTRITUBE_DIR = "./attributes"

type AttrKey struct {
	vendorId uint32
	attrType int
}

var (
	ATTRITUBES = map[AttrKey]*Attribute{}
)

// 属性值类型
type AttributeValueType int

const (
	INTEGE AttributeValueType = iota
	STRING
	OCTETS
	IP_ADDR
	OTHER
)

func getAttrValue(attributeType AttributeValueType, value []byte) string {
	var ret string
	switch attributeType {
	case INTEGE:
		val := binary.BigEndian.Uint32(value)
		ret = strconv.FormatUint(uint64(val), 10)
	case STRING:
		ret = string(value)
	case OCTETS:
		ret = strings.ToUpper(hex.EncodeToString(value))
	case IP_ADDR:
		return IPString(value)
	}
	return ret
}

func IPString(source []byte) string {
	if len(source) != 4 {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d.%d", source[0]&0xFF, source[1]&0xFF, source[2]&0xFF, source[3]&0xFF)
}

// 属性结构体
type Attribute struct {
	VendorId   uint32
	VendorName string
	Type       int
	Name       string
	// 枚举
	ValueType AttributeValueType
	// 属性值，可选的属性值
	AttributeValues []AttributeValue
}

// 属性值结构体
type AttributeValue struct {
	Name      string
	ValueName string
	Value     int
}

func ReadAttributeFiles() {
	infos, err := ioutil.ReadDir("attributes")
	if err != nil {
		log.Panicf("获取RADIUS属性文件列表失败，请检查文件目录是否正确, 错误：%s", err.Error())
	}
	for _, fileInfo := range infos {
		parseAttributes(fileInfo)
	}
}

// 解析属性文件
func parseAttributes(file os.FileInfo) {
	filePath := ATTRITUBE_DIR + "/" + file.Name()
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panicf("读取文件失败，文件名：%s, 错误：%s", filePath, err.Error())
	}

	pattern, _ := regexp.Compile("^#")
	splitPattern, _ := regexp.Compile("\\s+")
	buffer := bytes.NewBuffer(bs)

	var vendorId uint32 = 1
	vendorName := ""

	_attrNameType := make(map[string]int)

	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panicf("解析文件失败，请检查文件格式,file = %s, 错误 %s", filePath, err.Error())
		}

		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")
		if pattern.MatchString(line) || line == "" {
			continue
		}

		items := splitPattern.Split(line, -1)
		lineFirstItem := items[0]
		if len(items) == 3 && lineFirstItem == "VENDOR" {
			vendorName = items[1]
			val, err := strconv.Atoi(items[2])
			if err != nil {
				log.Panic(err)
			}
			vendorId = uint32(val)
		} else if len(items) >= 4 && lineFirstItem == "ATTRIBUTE" {
			typeVal, e := strconv.Atoi(items[2])
			typeName := items[1]
			valueType := getAttributeValueType(items[3])
			if e != nil {
				panic(e)
			}

			attr := Attribute{
				VendorId:   vendorId,
				VendorName: vendorName,
				Type:       typeVal,
				Name:       typeName,
				ValueType:  valueType,
			}
			ATTRITUBES[AttrKey{vendorId: vendorId, attrType: typeVal}] = &attr

			_attrNameType[typeName] = typeVal

		} else if len(items) == 4 && lineFirstItem == "VALUE" {
			belongAttrName := items[1]
			typeVal, ok := _attrNameType[belongAttrName]
			attribute, attrOk := ATTRITUBES[AttrKey{vendorId: vendorId, attrType: typeVal}]

			if !ok || !attrOk {
				continue
			}

			val, err := strconv.Atoi(items[3])
			if err != nil {
				continue
			}

			attrVal := AttributeValue{
				Name:      belongAttrName,
				ValueName: items[1],
				Value:     val,
			}
			values := attribute.AttributeValues
			if values == nil {
				values = make([]AttributeValue, 0, 20)
			}
			values = append(values, attrVal)
			attribute.AttributeValues = values
		}
	}
}

func getAttributeValueType(valueTypeName string) AttributeValueType {
	switch valueTypeName {
	case "string":
		return STRING
	case "integer":
		return INTEGE
	case "octets":
		return OCTETS
	case "ipaddr":
		return IP_ADDR
	default:
		return OTHER
	}
}
