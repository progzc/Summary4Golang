package go_explore

import (
	"fmt"
	"testing"
)

type People struct{}

func TestStruct_1(t *testing.T) {
	a := &People{}
	b := &People{}
	fmt.Println(a, b, a == b) // &{} &{} true
}

func TestStruct_2(t *testing.T) {
	a := &People{}
	b := &People{}
	fmt.Printf("%p\n", a)     // 0x696a58
	fmt.Printf("%p\n", b)     // 0x696a58
	fmt.Println(a, b, a == b) // &{} &{} true
}

// go run -gcflags="-m -l": 可以分析内存逃逸
// go run -gcflags="-N -l"： 可以禁用内存逃逸
func TestStruct_3(t *testing.T) {
	a := new(struct{})
	b := new(struct{})
	println(a, b, a == b) // 0xc00004bf4f 0xc00004bf4f false

	c := new(struct{})
	d := new(struct{})
	fmt.Println(c, d)     // &{} &{}
	println(c, d, c == d) // 0x697a78 0x697a78 true
}
