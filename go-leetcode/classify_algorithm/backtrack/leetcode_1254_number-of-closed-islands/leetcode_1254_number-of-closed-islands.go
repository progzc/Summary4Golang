package leetcode_1254_number_of_closed_islands

// 1254. 统计封闭岛屿的数目
// https://leetcode.cn/problems/number-of-closed-islands/

// closedIsland dfs
// 时间复杂度: O()
// 空间复杂度: O()
func closedIsland(grid [][]int) int {
	var (
		ans int
		dfs func(i, j int) bool
	)
	if len(grid) == 0 || len(grid[0]) == 0 {
		return ans
	}
	m, n := len(grid), len(grid[0])
	dfs = func(i, j int) bool {
		if i < 0 || i >= m || j < 0 || j >= n {
			return false
		}
		if grid[i][j] == 0 {
			grid[i][j] = 1
			// 下面这种是错误的：因为前面如果false后面就不执行了,没有把相邻的标记
			// return dfs(i-1, j) && dfs(i+1, j) && dfs(i, j-1) && dfs(i, j+1)
			up, down, left, right := dfs(i-1, j), dfs(i+1, j), dfs(i, j-1), dfs(i, j+1)
			return up && down && left && right

		}
		return true
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 0 {
				if dfs(i, j) {
					ans++
				}
			}
		}
	}
	return ans
}
