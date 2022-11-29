package leetcode_0221_maximal_square

// 221. 最大正方形
// https://leetcode.cn/problems/maximal-square/

// maximalSquare 动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路:
//	状态: dp[i][j] 表示以 (i, j) 为右下角，且只包含 1 的正方形的边长最大值。
//	转移方程: dp[i][j] = min{dp[i-1][j], dp[i-1][j-1], dp[i][j-1]}+1
func maximalSquare(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}

	maxSize := 0
	m, n := len(matrix), len(matrix[0])
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == '1' {
				if i == 0 || j == 0 {
					dp[i][j] = 1
				} else {
					dp[i][j] = min(min(dp[i-1][j], dp[i][j-1]), dp[i-1][j-1]) + 1
				}
				maxSize = max(maxSize, dp[i][j])
			}
		}
	}
	return maxSize * maxSize
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
