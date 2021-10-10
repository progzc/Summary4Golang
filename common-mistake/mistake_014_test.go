package common_mistake

import (
	"fmt"
	"testing"
)

// slice 和 array 其实是一维数据
func TestMistake_014(t *testing.T) {
	right014_1()
	right014_2()
}

// 使用原始的一维数组
// (1) 对每个内部 slice 进行内存分配
//	   注意：内部的 slice 相互独立，使得任一内部 slice 增缩都不会影响到其他的 slice
func right014_1() {
	x := 2
	y := 4

	table := make([][]int, x)
	for i := range table {
		table[i] = make([]int, y)
	}
}

// 使用“共享底层数组”的切片
// (1) 创建一个存放原始数据的容器 slice
// (2) 创建其他的 slice
// (3) 切割原始 slice 来初始化其他的 slice
func right014_2() {
	h, w := 2, 4
	raw := make([]int, h*w)

	for i := range raw {
		raw[i] = i
	}
	// 初始化原始 slice
	fmt.Println(raw, &raw[4]) // [0 1 2 3 4 5 6 7] 0xc420012120

	table := make([][]int, h)
	for i := range table {

		// 等间距切割原始 slice，创建动态多维数组 table
		// 0: raw[0*4: 0*4 + 4]
		// 1: raw[1*4: 1*4 + 4]
		table[i] = raw[i*w : i*w+w]
	}
	fmt.Println(table, &table[1][0]) // [[0 1 2 3] [4 5 6 7]] 0xc420012120
}
