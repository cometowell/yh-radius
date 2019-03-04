package main

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestLeftPadChar(t *testing.T) {
	fmt.Println(LeftPadChar("a", 'C', 20))
}

func TestRightPadChar(t *testing.T) {
	fmt.Println(RightPadChar("a", 'C', 20))
}

func TestBinary(t *testing.T) {
	container := make([]byte, 8)
	binary.BigEndian.PutUint64(container, 98)
	fmt.Println(container)
	binary.BigEndian.PutUint64(container, 789)
	fmt.Println(container)
}

func TestFillBytes(t *testing.T) {
	fmt.Println(FillBytesByString(64, "abasdeas"))
}