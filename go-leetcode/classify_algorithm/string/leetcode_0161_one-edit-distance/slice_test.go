package leetcode_0161_one_edit_distance

import (
	"fmt"
	"testing"
)

// 数组切片的下标范围：
// https://blog.csdn.net/yzf279533105/article/details/94745134
func TestSlice(t *testing.T) {
	input := []int{1, 2}
	// 这里并不会发生索引越界
	fmt.Println(input[2:])
	// 下面这里会发生索引越界
	//fmt.Println(input[3:])
}
