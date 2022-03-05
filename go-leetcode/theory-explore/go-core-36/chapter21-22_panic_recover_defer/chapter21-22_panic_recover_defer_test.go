package chapter21_22_panic_recover_defer

import (
	"fmt"
	"testing"
)

// TestPanic_1
// (1)Q:从 panic 被引发到程序终止运行的大致过程是什么?
//	   A:控制权从所属函数的那行代码一级一级地沿着调用栈的反方向传播至顶端，也就是我们编写的最外层函数那里
//	  	 (一般是main函数或者recover函数所在的goroutine),被Go runtime收回,最终程序崩溃并终止运行
// (2)Q:一个函数怎样才能把 panic 转化为error类型值，并将其作为函数的结果值返回给调用方?
//	  A: 使用defer+recover+debug.stack()

//func doSomething() (err error) {
//	defer func() {
//		p := recover()
//		err = fmt.Errorf("FATAL ERROR: %s", p)
//	}()
//	panic("Oops!!")
//}

// (3)Q:怎样让 panic 包含一个值，以及应该让它包含什么样的值?
//    A:最好传入 字符串 或 error
// (4)Q:怎样施加应对 panic 的保护措施，从而避免程序崩溃?
//    A:使用defer+recover
// (5)Q:如果一个函数中有多条defer语句，那么那几个defer函数调用的执行顺序是怎样的?
//	  A:先进后出，即同一个函数中，defer函数调用的执行顺序与它们分别所属的defer语句的出现顺序（更严谨地说，是执行顺序）完全相反

// (6)Q:我们可以在defer函数中恢复 panic，那么可以在其中引发 panic吗?
//	  A:当然可以。这样做可以把原先的 panic 包装一下再抛出去。
func TestPanic_1(t *testing.T) {
	//panic的传播（从里到外）
	fmt.Println("Enter function main.")
	caller1()
	fmt.Println("Exit function main.")
}

// TestPanic_2 defer语句的执行顺序
func TestPanic_2(t *testing.T) {
	defer fmt.Println("first defer")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("defer in for [%d]\n", i)
	}
	defer fmt.Println("last defer")
}

func caller1() {
	fmt.Println("Enter function caller1.")
	caller2()
	fmt.Println("Exit function caller1.")
}

func caller2() {
	fmt.Println("Enter function caller2.")
	s1 := []int{0, 1, 2, 3, 4}
	e5 := s1[5]
	_ = e5
	fmt.Println("Exit function caller2.")
}
