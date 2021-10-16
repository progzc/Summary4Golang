package closure

import "fmt"

func do1() {
	a := 1
	// 闭包,a会进行引用传递
	func() {
		fmt.Println(a)
		a = 2
	}()
}

// do1等价于do2
func do2() {
	a := 1
	func1(&a)
}

func func1(a *int) {
	fmt.Println(*a)
	*a = 2
}
