package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
)

// PAP 密码验证
func pap(shareSecret, password string, rp RadiusPackage) bool {
	hash := md5.New()
	hash.Write([]byte(shareSecret))
	hash.Write(rp.Authenticator[:])
	b := hash.Sum(nil)

	value := []byte(password)
	length := len(value)

	shortOf := length % 16
	times := length/16 + ifVal(shortOf > 0, 1, 0).(int)
	supplement := make([]byte, shortOf)
	value = append(value, supplement...)
	result := make([]byte, 0)
	for i := 0; i < times; i++ {
		left := value[i*16 : (i+1)*16]
		ret := [16]byte{}
		for j := 0; j < 16; j++ {
			ret[j] = left[j] ^ b[j]
		}
		result = append(result, ret[:]...)
	}

	attr := rp.RadiusAttrMap[AttrKey{Standard, 2}]
	return bytes.Equal(result, attr.AttrValue)
}

// CHAP认证 MD5(ID + PASSWORD明文 + CHALLENGE)
func chap(password string, rp *RadiusPackage) bool {
	var chapPassword []byte
	for _, attr := range rp.RadiusAttrs {
		if attr.AttrType == 3 {
			chapPassword = attr.AttrValue
		}
	}

	if len(chapPassword) != 17 {
		return false
	}

	var chapId = chapPassword[0]
	challenge := rp.challenge
	if len(challenge) != 16 {
		challenge = rp.Authenticator[:]
	}

	hashPassword := chapPassword[1:]

	buffer := bytes.NewBuffer(nil)
	buffer.WriteByte(chapId)
	buffer.Write([]byte(password))
	buffer.Write(challenge)
	sum := md5.Sum(buffer.Bytes())
	return bytes.Equal(sum[:], hashPassword)
}

func ifVal(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// 加密
func encrypt(origin string) string {
	key := config["encrypt.key"].(string)
	return AesEncrypt(origin, key)
}

// 解密
func decrypt(encryptedMsg string) string {
	key := config["encrypt.key"].(string)
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
