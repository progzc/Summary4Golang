package common_mistake

import (
	"fmt"
	"testing"
)

// Array 类型的值作为函数参数
// 在 Go 中，数组是值。作为参数传进函数时，传递的是数组的原始值拷贝，此时在函数内部是无法更新该数组的
func TestMistake_012(t *testing.T) {
	wrong012()
	right012_1()
	right012_2()
}

func wrong012() {
	x := [3]int{1, 2, 3}

	func(arr [3]int) {
		arr[0] = 7
		fmt.Println(arr) // [7 2 3]
	}(x)
	fmt.Println(x) // [1 2 3]    // 并不是你以为的 [7 2 3]
}

// 传址会修改原数据
func right012_1() {
	x := [3]int{1, 2, 3}
	func(arr *[3]int) {
		(*arr)[0] = 7
		fmt.Println(arr) // &[7 2 3]
	}(&x)
	fmt.Println(x) // [7 2 3]
}

// 直接使用 slice
func right012_2() {
	x := []int{1, 2, 3}
	func(arr []int) {
		arr[0] = 7
		fmt.Println(x) // [7 2 3]
	}(x)
	fmt.Println(x) // [7 2 3]
}
