package radius

import (
	"fmt"
	"go-rad/common"
	"testing"
)

func TestAesPassword(t *testing.T)  {
	aesEncrypt := common.AesEncrypt("123456", "vbRIKz90HJ$jjwyzu3BsUdci1600l7rP")
	fmt.Println(aesEncrypt)
}
