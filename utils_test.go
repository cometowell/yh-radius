package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
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

// 通过反射，对user进行赋值
type user struct {
	name    string
	age     int
	feature map[string]interface{}
}

func (u *user) test1() {

}

func TestStructReflect(t *testing.T) {
	var u interface{}
	u = new(user)
	value := reflect.ValueOf(u)
	if value.Kind() == reflect.Ptr {
		elem := value.Elem()
		name := elem.FieldByName("name")
		if name.Kind() == reflect.String {
			*(*string)(unsafe.Pointer(name.Addr().Pointer())) = "fangwendong"
		}

		age := elem.FieldByName("age")
		if age.Kind() == reflect.Int {
			*(*int)(unsafe.Pointer(age.Addr().Pointer())) = 24
		}

		feature := elem.FieldByName("feature")
		if feature.Kind() == reflect.Map {
			*(*map[string]interface{})(unsafe.Pointer(feature.Addr().Pointer())) = map[string]interface{}{
				"爱好": "篮球",
				"体重": 60,
				"视力": 5.2,
			}
		}

	}

	fmt.Println(u)
}
