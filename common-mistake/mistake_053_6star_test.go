package common_mistake

import (
	"fmt"
	"testing"
)

// 更新 map 字段的值
// (1) 如果 map 一个字段的值是 struct 类型，则无法直接更新该 struct 的单个字段
//     更新 map 中 struct 元素的字段值，有 2 个方法：
//     a. 使用局部变量
//     b. 使用指向元素的 map 指针
// (2) 需区分开的是，slice 的元素可寻址
func TestMistake_053(t *testing.T) {
	wrong053_1()
	right053_1()
}

// (1) 如果 map 一个字段的值是 struct 类型，则无法直接更新该 struct 的单个字段
func wrong053_1() {
	//type data struct {
	//	name string
	//}
	//m := map[string]data{
	//	"x": {"Tom"},
	//}
	//m["x"].name = "Jerry"
}

func right053_1() {
	type data struct {
		name string
	}
	m := map[string]data{
		"x": {"Tom"},
	}
	// a. 使用局部变量
	r := m["x"]
	r.name = "Jerry"
	m["x"] = r
	fmt.Println(m) // map[x:{Jerry}]
}

func right053_2() {
	type data struct {
		name string
	}
	// b. 使用指向元素的 map 指针
	m := map[string]*data{
		"x": {"Tom"},
	}

	m["x"].name = "Jerry" // 直接修改 m["x"] 中的字段
	fmt.Println(m["x"])   // &{Jerry}

	//--------注意下面这种误用-----------
	//m := map[string]*data{
	//	"x": {"Tom"},
	//}
	//m["z"].name = "what???"
	//fmt.Println(m["x"])
}

// (2) 需区分开的是，slice 的元素可寻址
func right053_3() {
	type data struct {
		name string
	}
	s := []data{{"Tom"}}
	s[0].name = "Jerry"
	fmt.Println(s) // [{Jerry}]
}
