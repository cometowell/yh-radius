package radius

import (
	"fmt"
	"go-rad/common"
	"testing"
)

func TestGetVlan(t *testing.T) {
	fmt.Println(common.getVlanIds(0,"vlanid=1;vlanid2=6"))
}

func TestUnknowLengthParam(t *testing.T) {
	fs := make([]f, 0)
	fs = append(fs, nil)
	fmt.Println(fs)
}

type f func()
