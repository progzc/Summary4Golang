package chapter15_pointer

import (
	"fmt"
	"testing"
)

// TestFuncMethod_1
// 关于指针方法和值方法: https://blog.csdn.net/weixin_44676081/article/details/111309791
// (1) 值对象可以调用值方法和指针方法
// (2) 指针对象可以调用值方法和指针方法
// (3) 若调用一个接口里面的函数,结构体对象实现接口时的方法可能是指针方法也可以是值方法,那么需要注意:
//  a.值类型只能调用值方法
//  b.指针类型可以调用值方法和指针方法
func TestFuncMethod_1(t *testing.T) {
	d1 := Dog{}
	d1.call()
	d1.phone() // 这里是语法糖, 相当于(&d1).phone()
	fmt.Println()

	d2 := &Dog{}
	d2.call() // 这里是语法糖, 相当于(*d1).call
	d2.phone()

	var jack Human
	jack = &Person{}
	jack.SayHello()
	jack.Cry()
	jack.Smile()

	//var tom Human
	//tom = Person{} // 值类型并没有实现Cry方法, 这里会编译报错
}

// -----------------------------------------------------
type Dog struct {
}

func (d Dog) call() {
	fmt.Println("call")
}

func (d *Dog) phone() {
	fmt.Println("phone")
}

// -----------------------------------------------------
type Human interface {
	SayHello()
	Cry()
	Smile()
}

type Person struct {
}

func (p Person) SayHello() {
	fmt.Println("SayHello")
}

func (p *Person) Cry() {
	fmt.Println("Cry")
}

func (p Person) Smile() {
	fmt.Println("Smile")
}
