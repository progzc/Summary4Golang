package common_mistake

import (
	"fmt"
	"testing"
)

// 使用指针作为方法的 receiver
// (1) 只要值是可寻址的，就可以在值上直接调用指针方法。即是对一个方法，它的 receiver 是指针就足矣。
// (2) 不是所有值都是可寻址的，比如 map 类型的元素、通过 interface 引用的变量
func TestMistake_052(t *testing.T) {
	wrong052()
	right052()
}

type data052 struct {
	name string
}

type printer interface {
	print()
}

func (p *data052) print() {
	fmt.Println("name: ", p.name)
}

func wrong052() {
	//d1 := data052{"one"}
	//d1.print() // d1 变量可寻址，可直接调用指针 receiver 的方法
	//
	//var in printer = data052{"two"}
	//in.print() // 类型不匹配
	//
	//m := map[string]data052{
	//	"x": data052{"three"},
	//}
	//m["x"].print() // m["x"] 是不可寻址的
}

func right052() {
	d1 := data052{"one"}
	d1.print() // d1 变量可寻址，可直接调用指针 receiver 的方法

	var in printer = &data052{"two"}
	in.print() // 类型不匹配

	m := map[string]*data052{
		"x": &data052{"three"},
	}
	m["x"].print() // m["x"] 是不可寻址的
}
