package main

import (
	"fmt"
	"testing"
)

func TestGetVlan(t *testing.T) {
	fmt.Println(standardGetVlanIds("vlanid=1;vlanid2=6"))
}
