package leetcode_0094_binary_tree_inorder_traversal

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	var result []int
	add(&result, 1)
	add(&result, 2)
	fmt.Println(result, len(result))
}

func add(result *[]int, num int) {
	*result = append(*result, num)
}

func TestSlice2(t *testing.T) {
	var result []int
	add2(result, 1)
	add2(result, 2)
	fmt.Println(result, len(result))
}

func add2(result []int, num int) {
	result = append(result, num)
}
