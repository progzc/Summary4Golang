package leetcode_0064_minimum_path_sum

// 0064.最小路径和
// https://leetcode-cn.com/problems/minimum-path-sum/

// minPathSum 动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(1)
// 思路：
//	a.递归表达式：f(i,j)=min{f(i,j-1),f(i-1,j)}
//	b.可以利用原数组进行数据存储，这样可以减少空间复杂度
func minPathSum(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	for i := 1; i < n; i++ {
		grid[0][i] = grid[0][i] + grid[0][i-1]
	}
	for i := 1; i < m; i++ {
		grid[i][0] = grid[i][0] + grid[i-1][0]
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			grid[i][j] += min(grid[i-1][j], grid[i][j-1])
		}
	}
	return grid[m-1][n-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
