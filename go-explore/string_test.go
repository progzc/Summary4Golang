package go_explore

import (
	"fmt"
	"testing"
)

// 参考文章 https://mp.weixin.qq.com/s/EbxkBokYBajkCR-MazL0ZA

type Student1 struct {
	Name string
	Age  int
}

type Student2 struct {
	Name string
	Age  int
}

type Student3 struct {
	Name string
	Age  int
}

func (s Student2) String() string {
	return fmt.Sprintf("[Name: %s], [Age: %d]", s.Name, s.Age)
}

func (s *Student3) String() string {
	return fmt.Sprintf("[Name: %s], [Age: %d]", s.Name, s.Age)
}

func TestString_1(t *testing.T) {
	var s = Student1{
		Name: "qcrao",
		Age:  18,
	}
	fmt.Println(s) // {qcrao 18}
}

func TestString_2(t *testing.T) {
	var s = Student2{
		Name: "qcrao",
		Age:  18,
	}
	// 类型 T 只有接受者是 T 的方法；而类型 *T 拥有接受者是 T 和 *T 的方法。语法上 T 能直接调 *T 的方法仅仅是 Go 的语法糖。
	fmt.Println(s)  // [Name: qcrao], [Age: 18]
	fmt.Println(&s) // [Name: qcrao], [Age: 18]
}

func TestString_3(t *testing.T) {
	var s = Student3{
		Name: "qcrao",
		Age:  18,
	}
	// 类型 T 只有接受者是 T 的方法；而类型 *T 拥有接受者是 T 和 *T 的方法。语法上 T 能直接调 *T 的方法仅仅是 Go 的语法糖。
	fmt.Println(s)  // {qcrao 18}
	fmt.Println(&s) // [Name: qcrao], [Age: 18]
}
