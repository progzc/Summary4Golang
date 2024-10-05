package leetcode_0074_search_a_2d_matrix

// 0074.搜索二维矩阵
// https://leetcode.cn/problems/search-a-2d-matrix

// searchMatrix
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(1)
func searchMatrix(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	for i := 0; i < m; i++ {
		if matrix[i][n-1] == target {
			return true
		} else if matrix[i][n-1] < target {
			continue
		} else {
			left, right := 0, n-1
			for left <= right {
				mid := left + (right-left)/2
				if matrix[i][mid] == target {
					return true
				} else if matrix[i][mid] < target {
					left = mid + 1
				} else {
					right = mid - 1
				}
			}
		}
	}
	return false
}
