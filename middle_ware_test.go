package main

import (
	"fmt"
	"testing"
)

func TestGetVlan(t *testing.T) {
	fmt.Println(getVlanIds("vlanid=1;vlanid2=6"))
}
