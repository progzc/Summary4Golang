package leetcode_0200_number_of_islands

// 0200.岛屿数量
// https://leetcode-cn.com/problems/number-of-islands/

// numIslands 深度优先遍历
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)，最坏情况下，政改革网格均为陆地
// 思路：在深度优先搜索的过程中，将每个搜索到的1都重新标记为0
func numIslands(grid [][]byte) int {
	if grid == nil || len(grid) == 0 {
		return 0
	}
	row, column := len(grid), len(grid[0])
	islands := 0
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			if grid[i][j] == '1' {
				islands++
				dfs(grid, i, j)
			}
		}
	}
	return islands
}

func dfs(grid [][]byte, row, column int) {
	if row < 0 || column < 0 || row >= len(grid) || column >= len(grid[0]) || grid[row][column] == '0' {
		return
	}
	grid[row][column] = '0'
	dfs(grid, row-1, column)
	dfs(grid, row+1, column)
	dfs(grid, row, column-1)
	dfs(grid, row, column+1)
}
