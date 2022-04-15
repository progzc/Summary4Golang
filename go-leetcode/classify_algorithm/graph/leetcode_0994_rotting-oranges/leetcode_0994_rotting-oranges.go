package leetcode_0994_rotting_oranges

// 0994.腐烂的橘子
// https://leetcode-cn.com/problems/rotting-oranges/

// orangesRotting 广度优先搜索
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：与二叉树的层序遍历有同样的思想
func orangesRotting(grid [][]int) int {
	// 定义坐标点及方向矢量
	type point struct {
		x, y int
	}
	dirs := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 上下左右

	// 初始化，寻找腐烂的橘子
	m, n := len(grid), len(grid[0])
	var stack []point
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 2 {
				stack = append(stack, point{i, j})
			}
		}
	}

	// 广度优先搜索
	count := 0
	for {
		find := false
		size := len(stack)
		for i := 0; i < size; i++ {
			p := stack[0]
			stack = stack[1:]
			for _, dir := range dirs {
				newX, newY := p.x+dir[0], p.y+dir[1]
				if newX >= 0 && newX < m && newY >= 0 && newY < n && grid[newX][newY] == 1 {
					find = true
					grid[newX][newY] = 2
					stack = append(stack, point{newX, newY})
				}
			}
		}
		if find == true {
			count++
		} else {
			break
		}
	}

	// 检查是否还有橘子未腐烂
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				return -1
			}
		}
	}
	return count
}
