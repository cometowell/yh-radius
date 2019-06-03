package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

func LeftPadChar(source string, padChar byte, size int) string {
	sourceLength := len(source)
	if sourceLength >= size {
		return source
	}
	return strings.Repeat(string(padChar), size-sourceLength) + source
}

func RightPadChar(source string, padChar byte, size int) string {
	sourceLength := len(source)
	if sourceLength >= size {
		return source
	}
	return source + strings.Repeat(string(padChar), size-sourceLength)
}

func IpAddrToBytes(ipAddr string) (ipArr []byte, err error) {
	ipArr = make([]byte, 4)
	items := strings.Split(ipAddr, ".")
	if len(items) != 4 {
		return nil, errors.New("ip地址格式错误")
	}

	for index, item := range items {
		val, e := strconv.Atoi(item)
		if e != nil {
			return nil, e
		}
		ipArr[index] = byte(val)
	}

	return ipArr, nil
}

func FillBytesByString(size int, value string) []byte {
	if len(value) >= size {
		return []byte(value)
	}
	ret := make([]byte, size)
	copy(ret, []byte(value))
	return ret
}

func GetIntegerBytes(val uint32) []byte {
	container := make([]byte, 4)
	binary.BigEndian.PutUint32(container, val)
	return container
}

// 加密
func Encrypt(origin string) string {
	key := GetConfig()["encrypt.key"].(string)
	return AesEncrypt(origin, key)
}

// 解密
func Decrypt(encryptedMsg string) string {
	key := GetConfig()["encrypt.key"].(string)
	return AesDecrypt(encryptedMsg, key)
}

func AesEncrypt(orig string, key string) string {
	origData := []byte(orig)
	k := []byte(key)
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	encryptedBytes := make([]byte, len(origData))
	blockMode.CryptBlocks(encryptedBytes, origData)
	return base64.StdEncoding.EncodeToString(encryptedBytes)
}
func AesDecrypt(encryptedMsg string, key string) string {
	encryptedBytes, _ := base64.StdEncoding.DecodeString(encryptedMsg)
	k := []byte(key)
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	orig := make([]byte, len(encryptedBytes))
	// 解密
	blockMode.CryptBlocks(orig, encryptedBytes)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// 补码
// AES加密数据块分组长度必须为128bit(byte[16])
// 密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPaddingLen := int(origData[length-1])
	return origData[:(length - unPaddingLen)]
}

