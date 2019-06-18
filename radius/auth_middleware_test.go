package radius

import (
	"fmt"
	"testing"
)

func TestGetVlan(t *testing.T) {
	fmt.Println(getVlanIds(0, "vlanid=1;vlanid2=6"))
}

func TestUnknowLengthParam(t *testing.T) {
	fs := make([]f, 0)
	fs = append(fs, nil)
	fmt.Println(fs)
}

type f func()
