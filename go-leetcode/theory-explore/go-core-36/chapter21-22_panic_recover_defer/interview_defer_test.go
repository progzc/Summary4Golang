package chapter21_22_panic_recover_defer

import (
	"errors"
	"fmt"
	"testing"
)

func Test_interview_defer_1(t *testing.T) {
	// 输出:
	// 3
	// 3
	// 3
	// 0
	defer1()
}

type number int

func (n number) print() {
	fmt.Println(n)
}

func (n *number) pprint() {
	fmt.Println(*n)
}

func defer1() {
	var n number
	defer n.print()
	defer n.pprint()
	defer func() {
		n.print()
	}()
	defer func() {
		n.pprint()
	}()
	n = 3
}

func Test_interview_defer_2(t *testing.T) {
	fmt.Println(defer2()) // 5
}

func defer2() int {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func Test_interview_defer_3(t *testing.T) {
	// 输出:
	// during return
	// before return
	defer3()
}

func defer3() {
	defer func() {
		fmt.Println("before return")
	}()

	if true {
		fmt.Println("during return")
		return
	}
	defer func() {
		fmt.Println("after return")
	}()
}

func Test_interview_defer_4(t *testing.T) {
	fmt.Println(defer4()) // 1
}

func defer4() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

func Test_interview_defer_5(t *testing.T) {
	fmt.Println(defer5()) // 6
}

func defer5() (r int) {
	defer func() {
		r = r + 5
	}()
	return 1
}

func Test_interview_defer_6(t *testing.T) {
	defer defer6()
}

func defer6() {
	f1() // <nil>
	f2() // defer2 error
	f3() // <nil>
}

func f1() {
	var err error
	defer fmt.Println(err)
	err = errors.New("defer1 error")
}

func f2() {
	var err error
	defer func() {
		fmt.Println(err)
	}()
	err = errors.New("defer2 error")
}

func f3() {
	var err error
	defer func(err error) {
		fmt.Println(err)
	}(err)
	err = errors.New("defer3 error")
}
