package common_mistake

import (
	"fmt"
	"testing"
	"time"
)

// for 语句中的迭代变量与闭包函数
// (1) for 语句中的迭代变量在每次迭代中都会重用，
//     即 for 中创建的闭包函数接收到的参数始终是同一个变量，在 goroutine 开始执行时都会得到同一个迭代值
func TestMistake_047(t *testing.T) {
	//wrong047()
	//right047_1()
	//right047_2()
	right047_3()
}

func wrong047() {
	data := []string{"one", "two", "three"}
	for _, v := range data {
		go func() {
			fmt.Println(v) // three three three
		}()
	}
	time.Sleep(3 * time.Second)
}

func right047_1() {
	data := []string{"one", "two", "three"}
	for _, v := range data {
		v := v
		go func() {
			fmt.Println(v) // one two three
		}()
	}
	time.Sleep(3 * time.Second)
}

func right047_2() {
	data := []string{"one", "two", "three"}
	for _, v := range data {
		go func(in string) {
			fmt.Println(in) // one two three
		}(v)
	}
	time.Sleep(3 * time.Second)
}

type field struct {
	name string
}

func (p *field) print() {
	fmt.Println(p.name)
}

func right047_3() {
	data := []*field{{"one"}, {"two"}, {"three"}}
	for _, v := range data { // 此时迭代值 v 是三个元素值的地址，每次 v 指向的值不同
		go v.print() // one two three
	}
	time.Sleep(3 * time.Second)
}
