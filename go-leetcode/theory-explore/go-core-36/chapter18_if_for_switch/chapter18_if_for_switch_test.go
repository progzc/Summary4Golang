package chapter18_if_for_switch

import (
	"fmt"
	"testing"
)

// TestFor_1 参透range的本质
// (1)Q:使用携带range子句的for语句时需要注意哪些细节?
//		range表达式的结果值可以是数组、数组的指针、切片、字符串、字典 或 接收通道 中的某一个，并且结果值只能有一个
//	  A:
//	 	(1)range表达式只会在for语句开始执行时被求值一次，无论后边会有多少次迭代
//		(2)range表达式的求值结果会被复制，也就是说，被迭代的对象是range表达式结果值的副本而不是原值
func TestFor_1(t *testing.T) {
	// 示例1: 切片的迭代
	numbers1 := []int{1, 2, 3, 4, 5, 6}
	for i := range numbers1 {
		if i == 3 {
			numbers1[i] |= i
		}
	}
	fmt.Println(numbers1) // [1 2 3 7 5 6]
	fmt.Println()

	// 示例2: 数组的迭代
	numbers2 := [...]int{1, 2, 3, 4, 5, 6}
	maxIndex2 := len(numbers2) - 1
	for i, e := range numbers2 { // 迭代的并不是numbers2,而是其副本;但由于数组是值类型的
		if i == maxIndex2 {
			numbers2[0] += e
		} else {
			numbers2[i+1] += e
		}
	}
	fmt.Println(numbers2) // [7 3 5 7 9 11]
	fmt.Println()

	// 示例3: 切片的迭代
	numbers3 := []int{1, 2, 3, 4, 5, 6}
	maxIndex3 := len(numbers3) - 1
	for i, e := range numbers3 { // 迭代的并不是numbers3,而是其副本;但由于切片是引用类型的
		if i == maxIndex3 {
			numbers3[0] += e
		} else {
			numbers3[i+1] += e // [22 3 6 10 15 21]
		}
	}
	fmt.Println(numbers3)
}

// TestSwitch_1
// (1)Q:switch语句中的switch表达式和case表达式之间有着怎样的联系?
//	  A:
//	 	a.case表达式一般由case关键字和一个表达式列表组成;表达式列表中的多个表达式之间需要有英文逗号,分割
//		b.只有类型相同的值之间才有可能被允许进行判等操作
//		c.如果case表达式中子表达式的结果值是无类型的常量，那么它的类型会被自动地转换为switch表达式的结果类型;反之则不能转换
// (2)Q:switch语句对它的case表达式有哪些约束?
//	  A:
//	   a.switch语句不允许case表达式中的子表达式结果值存在相等的情况(这条约束只针对 由字面量直接表示的子表达式)
//	   b.上面这种绕过方式对用于类型判断的switch语句（以下简称为类型switch语句）就无效了
// (3)Q:在类型switch语句中，我们怎样对被判断类型的那个值做相应的类型转换?
//	   a. switch t := x.(type) {case}
// (4)Q:在if语句中，初始化子句声明的变量的作用域是什么?
//	   a.如果这个变量是新的变量，那么它的作用域就是当前if语句所代表的代码块。
//	     注意，后续的else if子句和else子句也包含在当前的if语句代表的代码块之内
func TestSwitch_1(t *testing.T) {
	//c.如果case表达式中子表达式的结果值是无类型的常量，那么它的类型会被自动地转换为switch表达式的结果类型;反之则不能转换
	// 示例1: switch是无类型的常量,case是int8
	//value1 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	//switch 1 + 3 { // 这条语句无法编译通过:只有类型相同的值之间才有可能被允许进行判等操作
	//case value1[0], value1[1]:
	//	fmt.Println("0 or 1")
	//case value1[2], value1[3]:
	//	fmt.Println("2 or 3")
	//case value1[4], value1[5], value1[6]:
	//	fmt.Println("4 or 5 or 6")
	//}

	// 示例2: switch是int8,case是无类型的常量
	value2 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value2[4] {
	case 0, 1:
		fmt.Println("0 or 1")
	case 2, 3:
		fmt.Println("2 or 3")
	case 4, 5, 6:
		fmt.Println("4 or 5 or 6")
	}

	// ---------------------------------------------------------------------------------
	//a.switch语句不允许case表达式中的子表达式结果值存在相等的情况(这条约束只针对结果值为常量的子表达式)
	// 示例1。
	//value3 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	//switch value3[4] { // 这条语句无法编译通过。
	//case 0, 1, 2:
	//	fmt.Println("0 or 1 or 2")
	//case 2, 3, 4:
	//	fmt.Println("2 or 3 or 4")
	//case 4, 5, 6:
	//	fmt.Println("4 or 5 or 6")
	//}

	// 示例2。
	value5 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value5[4] {
	case value5[0], value5[1], value5[2]:
		fmt.Println("0 or 1 or 2")
	case value5[2], value5[3], value5[4]:
		fmt.Println("2 or 3 or 4")
	case value5[4], value5[5], value5[6]:
		fmt.Println("4 or 5 or26")
	}

	// ---------------------------------------------------------------------------------
	//b.上面这种绕过方式对用于类型判断的switch语句（以下简称为类型switch语句）就无效了
	// 示例3。
	//value6 := interface{}(byte(127))
	//switch t := value6.(type) { // 这条语句无法编译通过: byte类型是uint8的别名
	//case uint8, uint16:
	//	fmt.Println("uint8 or uint16")
	//case byte:
	//	fmt.Printf("byte")
	//default:
	//	fmt.Printf("unsupported type: %T", t)
	//}
}
