package main

import (
	"fmt"
	"testing"
)

func TestGetVlan(t *testing.T) {
	fmt.Println(getVlanIds(0,"vlanid=1;vlanid2=6"))
}
