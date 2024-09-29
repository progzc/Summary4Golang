package chapter12_func

import (
	"errors"
	"fmt"
	"testing"
)

// TestFunc_1 函数的使用
// (1)各个参数和结果的名称,甚至函数的名称都不能算作函数签名的一部分,甚至对于结果声明来说,没有名称都可以
// (2)高阶函数: 接受其他的函数作为参数传入;把其他的函数作为结果返回
// (3)函数类型属于引用类型，它的值可以为nil
// (4)闭包: 在一个函数中存在对外来自由变量的引用
// (5)函数传递的方式: 值传递,注意 切片/字典/通道 为引用类型, 【其中切片的引用指向的是层数组中某一个元素的指针】
// (6)函数真正拿到的参数值其实只是它们的副本,那么函数返回给调用方的结果值也会被复制
func TestFunc_1(t *testing.T) {
	// 函数的基本使用
	var p Printer
	p = printToStd
	p("something")

	// -----------------------------------------------------------------
	// (2)高阶函数: 接受其他的函数作为参数传入;把其他的函数作为结果返回
	// 方案一
	x, y := 12, 23
	op := func(x, y int) int {
		return x + y
	}
	result, err := calculate(x, y, op)
	fmt.Printf("The result: %d (error: %v)\n",
		result, err)
	result, err = calculate(x, y, nil)
	fmt.Printf("The result: %d (error: %v)\n",
		result, err)

	// 方案二
	x, y = 56, 78
	add := genCalculator(op)
	result, err = add(x, y)
	fmt.Printf("The result: %d (error: %v)\n",
		result, err)

	// -----------------------------------------------------------------
	// (5)函数传递的方式: 值传递
	array1 := [3]string{"a", "b", "c"}
	fmt.Printf("The array: %v\n", array1)
	array2 := modifyArray(array1)
	fmt.Printf("The modified array: %v\n", array2)
	fmt.Printf("The original array: %v\n", array1)
	fmt.Println()

	slice1 := []string{"x", "y", "z"}
	fmt.Printf("The slice: %v\n", slice1)
	slice2 := modifySlice(slice1)
	fmt.Printf("The modified slice: %v\n", slice2)
	fmt.Printf("The original slice: %v\n", slice1)
	fmt.Println()

	complexArray1 := [3][]string{
		[]string{"d", "e", "f"},
		[]string{"g", "h", "i"},
		[]string{"j", "k", "l"},
	}
	fmt.Printf("The complex array: %v\n", complexArray1)
	complexArray2 := modifyComplexArray(complexArray1)
	fmt.Printf("The modified complex array: %v\n", complexArray2)
	fmt.Printf("The original complex array: %v\n", complexArray1)
}

// ---------------------------------------------------------------
type Printer func(contents string) (n int, err error)

func printToStd(contents string) (bytesNum int, err error) {
	return fmt.Println(contents)
}

// -----------------------方案一----------------------------------
type operate func(x, y int) int

func calculate(x int, y int, op operate) (int, error) {
	if op == nil {
		return 0, errors.New("invalid operation")
	}
	return op(x, y), nil
}

// ------------------------方案二---------------------------------
type calculateFunc func(x int, y int) (int, error)

func genCalculator(op operate) calculateFunc {
	return func(x int, y int) (int, error) {
		// 这里op会出现闭包
		if op == nil {
			return 0, errors.New("invalid operation")
		}
		return op(x, y), nil
	}
}

// --------------------------------------------------------------
func modifyArray(a [3]string) [3]string {
	a[1] = "x"
	return a
}

func modifySlice(a []string) []string {
	a[1] = "i"
	return a
}

func modifyComplexArray(a [3][]string) [3][]string {
	a[1][1] = "s"
	a[2] = []string{"o", "p", "q"}
	return a
}
