package go_core_36

import (
	"errors"
	"fmt"
	"testing"
)

type Printer func(contents string) (n int, err error)

func printToStd(contents string) (bytesNum int, err error) {
	return fmt.Println(contents)
}

// TestFunc 函数声明与函数签名
func TestFunc(t *testing.T) {
	var p Printer
	p = printToStd
	p("something")
}

type operate func(x, y int) int

func calculate(x int, y int, op operate) (int, error) {
	if op == nil {
		return 0, errors.New("invalid operation")
	}
	return op(x, y), nil
}

// TestFunc2 高阶函数的使用
func TestFunc2(t *testing.T) {
	// 1.高阶函数的条件一: 接受其他的函数作为参数传入；
	// 2.高阶函数的条件二: 把其他的函数作为结果返回。
	// 只要满足其中一个条件即为高阶函数
	sum, _ := calculate(1, 2, func(x, y int) int {
		return x + y
	})
	fmt.Println(sum)
}

type calculateFunc func(x int, y int) (int, error)

func genCalculator(op operate) calculateFunc {
	// 使用高阶函数实现闭包
	// 实现闭包的意义: 动态地生成程序逻辑
	return func(x int, y int) (int, error) {
		if op == nil {
			return 0, errors.New("invalid operation")
		}
		return op(x, y), nil
	}
}

// TestFunc3 高阶函数的使用
func TestFunc3(t *testing.T) {
	// 1.高阶函数的条件一: 接受其他的函数作为参数传入；
	// 2.高阶函数的条件二: 把其他的函数作为结果返回。
	// 只要满足其中一个条件即为高阶函数
	x, y := 56, 78
	op := func(x, y int) int {
		return x + y
	}
	add := genCalculator(op)
	result, err := add(x, y)
	fmt.Printf("The result: %d (error: %v)\n",
		result, err)
}

// TestFunc4 值传递
func TestFunc4(t *testing.T) {
	array1 := [3]string{"a", "b", "c"}
	fmt.Printf("The array: %v\n", array1) // [a b c]
	array2 := modifyArray(array1)
	fmt.Printf("The modified array: %v\n", array2) // [a x c]
	fmt.Printf("The original array: %v\n", array1) // [a b c]
}
func modifyArray(a [3]string) [3]string {
	a[1] = "x"
	return a
}

// TestFunc5 引用传递
func TestFunc5(t *testing.T) {
	array1 := [3][]string{
		[]string{"d", "e", "f"},
		[]string{"g", "h", "i"},
		[]string{"j", "k", "l"},
	}
	fmt.Printf("The array: %v\n", array1) // [[d e f] [g h i] [j k l]]
	array2 := modifyArray2(array1)
	fmt.Printf("The modified array: %v\n", array2) // [[d e f] [x h i] [j k l]]
	fmt.Printf("The original array: %v\n", array1) // [[d e f] [x h i] [j k l]]
}

func modifyArray2(a [3][]string) [3][]string {
	a[1][0] = "x"
	return a
}
