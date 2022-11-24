package futu_go_basic

import (
	"errors"
	"fmt"
	"testing"
)

func Test_Slice(t *testing.T) {
	s := []int{1, 2, 3, 4}
	Append(s)
	fmt.Println(s) // [1 2 3 4]
	Add(s)
	fmt.Println(s) // [6 7 8 9]
}

func Append(s []int) {
	s = append(s, 5)
}

func Add(s []int) {
	for i := range s {
		s[i] = s[i] + 5
	}
}

type T interface{}

func Test_Interface(t *testing.T) {
	var (
		tt T
		p  *T
		i1 interface{} = tt
		i2 interface{} = p
	)
	fmt.Println(i1 == tt, i1 == nil) // true true
	fmt.Println(i2 == p, i2 == nil)  //  true false
}

func Test_Defer(t *testing.T) {
	fmt.Println(foo()) // a
}

func foo() (err error) {
	defer func() {
		fmt.Println("defer1", err) // defer1 c
		err = errors.New("a")
	}()
	// defer 语句表达式的值在定义时就已经确定了
	defer func(e error) {
		fmt.Println("defer2", e) // defer2 <nil>
		e = errors.New("b")
	}(err)
	return errors.New("c")
}
