package main

import (
	"bytes"
	"errors"
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
	attributes = map[AttrKey]Attribute{}
)

// 属性值类型
type AttributeValueType int

const (
	INTEGE AttributeValueType = iota
	STRING
	OCTETS
	IP_ADDR
)

// 属性结构体
type Attribute struct {
	VendorId uint32
	VendorName string
	Type int
	Name string
	// 枚举
	ValueType AttributeValueType
	// 属性值
	AttributeValue AttributeValue
}

// 属性值结构体
type AttributeValue struct {
	Name string
	ValueName string
	Value []int
}

func readAttributeFiles() {
	infos, err := ioutil.ReadDir("attributes")
	if err != nil {
		log.Panicf("获取RADIUS属性文件列表失败，请检查文件目录是否正确, 错误：%s", err.Error())
	}
	for _, fileInfo := range infos {
		parseAttributes(fileInfo)
		break
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

	var vendorId uint32 = 0
	vendorName := ""

	for  {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Panicf("解析文件失败，请检查文件格式,file = %s, 错误 %s", filePath, err.Error())
		}

		line = strings.Replace(line, "\n", "", -1)
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
		} else if len(items) == 4 && lineFirstItem == "ATTRIBUTE" {
			typeVal, e := strconv.Atoi(items[2])
			typeName := items[1]
			valueType, e := getAttributeValueType(items[3])
			if e != nil {
				panic(e)
			}

			attr := Attribute{
				VendorId: vendorId,
				VendorName: vendorName,
				Type: typeVal,
				Name: typeName,
				ValueType: valueType,
			}
			attributes[AttrKey{vendorId:vendorId, attrType:typeVal}] = attr

		} else if len(items) == 4 && lineFirstItem == "VALUE" {

		}
		break
	}
}

func getAttributeValueType(valueTypeName string) (AttributeValueType, error) {
	switch valueTypeName {
	case "string":
		return STRING, nil
	case "integer":
		return INTEGE, nil
	case "octets":
		return OCTETS, nil
	case "ipaddr":
		return IP_ADDR, nil
	}
	return 0, errors.New("找不到匹配类型")
}