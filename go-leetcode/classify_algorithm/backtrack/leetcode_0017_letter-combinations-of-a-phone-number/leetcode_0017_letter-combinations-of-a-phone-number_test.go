package leetcode_0017_letter_combinations_of_a_phone_number

import (
	"fmt"
	"testing"
)

// TestSlice
// 注意事项：如果有追加操作，则需要用切片指针；如果只是修改，则可以使用切片
func TestSlice(t *testing.T) {
	str1 := []string{"a", "b", "c"}
	str2 := []string{"e", "f"}
	result1 := appendSlice(str1)
	result2 := appendSlicePointer(&str2)
	fmt.Printf("str1: %v; result1： %v\n", str1, result1) // str1: [a b c]; result1： [a b c d]
	fmt.Printf("str2: %v; result2： %v\n", str2, result2) // str2: [e f g]; result2： [e f g]

}

func appendSlice(str []string) []string {
	str = append(str, "d")
	return str
}

func appendSlicePointer(str *[]string) []string {
	*str = append(*str, "g")
	return *str
}
