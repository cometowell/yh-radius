package web

import (
	"fmt"
	"testing"
)

func TestNilLen(t *testing.T) {
	var a []interface{} = nil
	fmt.Println(a)
	fmt.Println(len(a))
}
