package main

import (
	"fmt"
	"testing"
)

func TestAesPassword(t *testing.T)  {
	aesEncrypt := AesEncrypt("111111", "vbRIKz90HJ$jjwyzu3BsUdci1600l7rP")
	fmt.Println(aesEncrypt)
}
