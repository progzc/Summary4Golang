package leetcode_0659_max_area_of_island

// 0695. 岛屿的最大面积
// https://leetcode.cn/problems/max-area-of-island/

// maxAreaOfIsland 深度优先遍历
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
func maxAreaOfIsland(grid [][]int) int {
	if grid == nil || len(grid) == 0 {
		return 0
	}

	row, column := len(grid), len(grid[0])
	area := 0
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			if grid[i][j] == 1 {
				cnt := dfs(grid, i, j)
				if cnt > area {
					area = cnt
				}
			}
		}
	}
	return area
}

func dfs(grid [][]int, row, column int) int {
	var ans int
	if row < 0 || column < 0 || row >= len(grid) || column >= len(grid[0]) || grid[row][column] == 0 {
		return ans
	}
	grid[row][column] = 0
	ans++
	ans += dfs(grid, row-1, column)
	ans += dfs(grid, row+1, column)
	ans += dfs(grid, row, column-1)
	ans += dfs(grid, row, column+1)
	return ans
}
