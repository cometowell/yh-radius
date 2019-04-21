package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"
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

func TestMonthLastTime(t *testing.T) {
	fmt.Println(getMonthLastTime())
}

func TestDefaultTime(t *testing.T) {
	t1, _ := getStdTimeFromString("2099-12-31 23:59:59")
	fmt.Println(t1)
}

type A struct {
	B string
}

func TestSlice(t *testing.T) {
	l := []A{{"F"}, {"A"}, {"c"}, {"b"}, {"H"}}
	ls := make([]*A, 0)
	for _, item := range l {
		r := item
		//ls = append(ls, &item) // item指向同一内存地址，表明：item变量用于存储循环变量，地址不变
		ls = append(ls, &r)
	}
	fmt.Println(ls)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func TestReg(t *testing.T) {
	compile, _ := regexp.Compile("/a/b")
	fmt.Println(compile.Match([]byte("")))
}

func TestSlice2(t *testing.T) {
	var users []RadUser
	fmt.Println(len(users))
}
