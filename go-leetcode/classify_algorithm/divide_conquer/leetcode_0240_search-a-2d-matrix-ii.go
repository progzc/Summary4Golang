package divide_conquer

import "sort"

// 0240.搜索二维矩阵 II
// https://leetcode-cn.com/problems/search-a-2d-matrix-ii/

// searchMatrix 二分查找
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(1)
// 思路：对每行进行二分查找
func searchMatrix(matrix [][]int, target int) bool {
	for _, row := range matrix {
		idx := sort.SearchInts(row, target)
		// 注意：下面这行写成： if idx < len(row) {return true} 会报错
		// 因为：sort.SearchInts(row, target)，当搜索不到target时，返回的是插入的位置
		if idx < len(row) && row[idx] == target {
			return true
		}
	}
	return false
}

// searchMatrix_2 z字型查找
// 时间复杂度: O(m+n)
// 空间复杂度: O(1)
// 思路：从右上角开始搜索，每次可以忽略一行或一列
func searchMatrix_2(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	x, y := 0, n-1
	for x < m && y >= 0 {
		if matrix[x][y] == target {
			return true
		}
		if matrix[x][y] > target {
			y--
		} else {
			x++
		}
	}
	return false
}
