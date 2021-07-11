package go_core_36

import (
	"fmt"
	"testing"
)

// TestRange range表达式
func TestRange(t *testing.T) {
	// 示例1。
	numbers1 := []int{1, 2, 3, 4, 5, 6}
	// 若迭代变量只有一个，则代表索引值
	for i := range numbers1 {
		if i == 3 {
			numbers1[i] |= i
		}
	}
	fmt.Println(numbers1)
	fmt.Println()

	// 示例2：数组的迭代
	numbers2 := [...]int{1, 2, 3, 4, 5, 6}
	maxIndex2 := len(numbers2) - 1
	// 若迭代变量有多个，则左边的代表索引值，右边的代表元素值
	for i, e := range numbers2 {
		if i == maxIndex2 {
			numbers2[0] += e
		} else {
			numbers2[i+1] += e
		}
	}
	fmt.Println(numbers2)
	fmt.Println()

	// 示例3：切片的迭代
	numbers3 := []int{1, 2, 3, 4, 5, 6}
	maxIndex3 := len(numbers2) - 1
	for i, e := range numbers3 {
		if i == maxIndex3 {
			numbers3[0] += e
		} else {
			numbers3[i+1] += e
		}
	}
	fmt.Println(numbers3)
}

func TestSwitch(t *testing.T) {
	// 示例1。
	//value1 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	//switch 1 + 3 { // 这条语句无法编译通过。
	//case value1[0], value1[1]:
	//	fmt.Println("0 or 1")
	//case value1[2], value1[3]:
	//	fmt.Println("2 or 3")
	//case value1[4], value1[5], value1[6]:
	//	fmt.Println("4 or 5 or 6")
	//}

	// 示例2。
	value2 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value2[4] {
	// case中会自动转换成switch中的类型
	// 注意事项：以switch表达式为基准的自动类型转换仅在case子表达式的结果值为无类型常量时才会发生
	case 0, 1:
		fmt.Println("0 or 1")
	case 2, 3:
		fmt.Println("2 or 3")
	case 4, 5, 6:
		fmt.Println("4 or 5 or 6")
	}
}

func TestSwitch2(t *testing.T) {
	// 1. switch语句不允许case表达式中的子表达式结果值存在相等的情况，
	//    不论这些结果值相等的子表达式，是否存在于不同的case表达式中，都会是这样的结果。
	// 上面这个约束只针对结果值为常量的子表达式
	//value3 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	//switch value3[4] { // 这条语句无法编译通过。
	//case 0, 1, 2:
	//	fmt.Println("0 or 1 or 2")
	//case 2, 3, 4:
	//	fmt.Println("2 or 3 or 4")
	//case 4, 5, 6:
	//	fmt.Println("4 or 5 or 6")
	//}

	// 示例2：破解case相等的约束
	value5 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value5[2] {
	case value5[0], value5[1], value5[2]:
		fmt.Println("0 or 1 or 2")
	case value5[2], value5[3], value5[4]:
		fmt.Println("2 or 3 or 4")
	case value5[4], value5[5], value5[6]:
		fmt.Println("4 or 5 or26")
	}

	// 示例3：破解case相等的约束的方法不适用于类型判断
	// 因为类型switch语句中的case表达式的子表达式，都必须直接由类型字面量表示，而无法通过间接的方式表示。
	//value6 := interface{}(byte(127))
	//switch t := value6.(type) { // 这条语句无法编译通过。
	// byte时uint8的别名，本质上类型是一致的
	//case uint8, uint16:
	//	fmt.Println("uint8 or uint16")
	//case byte:
	//	fmt.Printf("byte")
	//default:
	//	fmt.Printf("unsupported type: %T", t)
	//}
}
