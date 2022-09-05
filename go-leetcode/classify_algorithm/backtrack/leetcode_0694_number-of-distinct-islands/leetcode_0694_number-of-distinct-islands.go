package leetcode_0694_number_of_distinct_islands

import (
	"fmt"
	"strings"
)

// 0694. 不同岛屿的数量
// https://leetcode.cn/problems/number-of-distinct-islands/

// numDistinctIslands 深度优先遍历
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)，最坏情况下，政改革网格均为陆地
// 思路：
//	a.在深度优先搜索的过程中，将每个搜索到的1都重新标记为0
//	b.采用 自定义序列化路径代号 + 散列表 来进行去重
func numDistinctIslands(grid [][]int) int {
	if grid == nil || len(grid) == 0 {
		return 0
	}
	row, column := len(grid), len(grid[0])
	m := make(map[string]int)
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			if grid[i][j] == 1 {
				sb := &strings.Builder{}
				dfs(grid, i, j, sb, 0)
				m[sb.String()] += 1
			}
		}
	}
	return len(m)
}

func dfs(grid [][]int, row, column int, sb *strings.Builder, dir int) {
	if row < 0 || column < 0 || row >= len(grid) || column >= len(grid[0]) || grid[row][column] == 0 {
		return
	}
	grid[row][column] = 0

	sb.WriteString(fmt.Sprintf("%d", dir))
	dfs(grid, row-1, column, sb, 1)
	dfs(grid, row+1, column, sb, 2)
	dfs(grid, row, column-1, sb, 3)
	dfs(grid, row, column+1, sb, 4)
	sb.WriteString(fmt.Sprintf("%d", -dir))
}
