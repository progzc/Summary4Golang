package common_mistake

import (
	"fmt"
	"testing"
)

// nil interface 和 nil interface 值
// (1) 虽然 interface 看起来像指针类型，但它不是; interface 类型的变量只有在类型和值均为 nil 时才为 nil
// (2) 如果你的 interface 变量的值是跟随其他变量变化的，与 nil 比较相等时小心
func TestMistake_054(t *testing.T) {
	wrong054()
	right054()
	right054_1()
}

func wrong054() {
	doIt := func(arg int) interface{} {
		var result *struct{} = nil
		if arg > 0 {
			result = &struct{}{}
		}
		return result
	}

	if res := doIt(-1); res != nil {
		fmt.Println("Good result: ", res) // Good result:  <nil>
		fmt.Printf("%T\n", res)           // *struct {}    // res 不是 nil，它的值为 nil
		fmt.Printf("%v\n", res)           // <nil>
	}
}

func right054() {
	doIt := func(arg int) interface{} {
		var result *struct{} = nil
		if arg > 0 {
			result = &struct{}{}
		} else {
			return nil // 明确指明返回 nil
		}
		return result
	}

	if res := doIt(-1); res != nil {
		fmt.Println("Good result: ", res)
	} else {
		fmt.Println("Bad result: ", res) // Bad result:  <nil>
	}
}

func right054_1() {
	var data *byte
	var in interface{}

	fmt.Println(data, data == nil) // <nil> true
	fmt.Println(in, in == nil)     // <nil> true

	in = data
	fmt.Println(in, in == nil) // <nil> false    // data 值为 nil，但 in 值不为 nil
}
