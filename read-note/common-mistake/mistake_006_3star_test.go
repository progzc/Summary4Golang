package common_mistake

import (
	"fmt"
	"testing"
)

type info struct {
	result int
}

func work() (int, error) {
	return 3, nil
}

// 不能使用简短声明来设置字段的值
// struct 的变量字段不能使用 := 来赋值以使用预定义的变量来避免解决：
func TestMistake_006(t *testing.T) {
	wrong006()
	right006()
}

func wrong006() {
	//var data info
	//data.result, err := work()    // error: non-name data.result on left side of :=
	//fmt.Printf("info: %+v\n", data)
}

func right006() {
	var data info
	var err error // err 需要预声明

	data.result, err = work()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("info: %+v\n", data)
}
