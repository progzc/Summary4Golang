package leetcode_0064_minimum_path_sum

// 0064.最小路径和
// https://leetcode-cn.com/problems/minimum-path-sum/

// minPathSum 动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(1)
// 思路：
//
//	a.递归表达式：f(i,j)=min{f(i,j-1),f(i-1,j)}
//	b.可以利用原数组进行数据存储，这样可以减少空间复杂度
//
// 条件：可以改变原数组
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

// minPathSum_2 动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(n)
// 思路：
//
//	a.递归表达式：f(i,j)=min{f(i,j-1),f(i-1,j)}
//
// 条件：不可以改变原数组
func minPathSum_2(grid [][]int) int {
	m, n := len(grid), len(grid[0])
	dp := make([]int, n)
	dp[0] = grid[0][0]

	for i := 1; i < n; i++ {
		dp[i] = grid[0][i] + dp[i-1]
	}

	for i := 1; i < m; i++ {
		dp[0] = dp[0] + grid[i][0]
		for j := 1; j < n; j++ {
			dp[j] = min(dp[j], dp[j-1]) + grid[i][j]
		}
	}
	return dp[n-1]
}

// minPathSum_3 动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：
//
//	a.递归表达式：f(i,j)=min{f(i,j-1),f(i-1,j)}
//
// 条件：不可以改变原数组
func minPathSum_3(grid [][]int) int {
	m := len(grid)
	n := len(grid[0])
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		if i == 0 {
			dp[i][0] = grid[i][0]
		} else {
			dp[i][0] = dp[i-1][0] + grid[i][0]
		}
	}

	for j := 0; j < n; j++ {
		if j == 0 {
			dp[0][j] = grid[0][j]
		} else {
			dp[0][j] = dp[0][j-1] + grid[0][j]
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			dp[i][j] = min(dp[i-1][j]+grid[i][j], dp[i][j-1]+grid[i][j])
		}
	}
	return dp[m-1][n-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
