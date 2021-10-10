package common_mistake

import (
	"fmt"
	"testing"
)

// 在 range 迭代 slice、array、map 时通过更新引用来更新元素
// (1) 在 range 迭代中，得到的值其实是元素的一份值拷贝，更新拷贝并不会更改原来的元素，即是拷贝的地址并不是原有元素的地址
func TestMistake_041(t *testing.T) {
	wrong041()
	right041_1()
}

func wrong041() {
}

func right041_1() {
	data := []int{1, 2, 3}
	for _, v := range data {
		v *= 10 // data 中原有元素是不会被修改的
	}
	fmt.Println("data: ", data) // data:  [1 2 3]
}

func right041_2() {
	data := []int{1, 2, 3}
	for i, v := range data {
		data[i] = v * 10 // 如果要修改原有元素的值，应该使用索引直接访问
	}
	fmt.Println("data: ", data) // data:  [10 20 30]
}

func right041_3() {
	data := []*struct{ num int }{{1}, {2}, {3}}
	for _, v := range data {
		v.num *= 10 // 直接使用指针更新
	}
	fmt.Println(data[0], data[1], data[2]) // &{10} &{20} &{30}
}
