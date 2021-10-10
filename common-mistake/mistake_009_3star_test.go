package common_mistake

import (
	"testing"
)

// 直接使用值为 nil 的 slice、map
// 允许对值为 nil 的 slice 添加元素，但对值为 nil 的 map 添加元素则会造成运行时 panic
func TestMistake_009(t *testing.T) {
	wrong009()
	right009()
}

func wrong009() {
	var m map[string]int
	m["one"] = 1 // error: panic: assignment to entry in nil map
	// m := make(map[string]int)// map 的正确声明，分配了实际的内存
}

func right009() {
	var s []int
	s = append(s, 1)
}
