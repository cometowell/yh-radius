package main

import (
	"fmt"
	"testing"
)

func TestGetVlan(t *testing.T) {
	fmt.Println(standardGetVlanIds("aaavlanid=1;vlanid2=6sdfasfasfasdfas"))
}
