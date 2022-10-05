package leetcode_0505_the_maze_ii

import "math"

// 0505. 迷宫 II
// https://leetcode.cn/problems/the-maze-ii/

// shortestDistance dfs
// 时间复杂度: O(m*n*max(m,n))
// 空间复杂度: O(m*n)
func shortestDistance(maze [][]int, start []int, destination []int) int {
	m, n := len(maze), len(maze[0])
	// distance[i][j]: 记录从起始位置到 (i, j) 的最小步数
	distance := make([][]int, m)
	for i := 0; i < m; i++ {
		distance[i] = make([]int, n)
		for j := 0; j < n; j++ {
			distance[i][j] = math.MaxInt32
		}
	}
	distance[start[0]][start[1]] = 0

	var (
		dfs  func(start []int)
		dirs = [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
	)
	dfs = func(start []int) {
		for _, dir := range dirs {
			x, y := start[0]+dir[0], start[1]+dir[1]
			count := 0
			for x >= 0 && y >= 0 && x < m && y < n && maze[x][y] == 0 {
				x += dir[0]
				y += dir[1]
				count++
			}

			if distance[start[0]][start[1]]+count < distance[x-dir[0]][y-dir[1]] {
				distance[x-dir[0]][y-dir[1]] = distance[start[0]][start[1]] + count
				dfs([]int{x - dir[0], y - dir[1]})
			}
		}
	}

	dfs(start)
	if distance[destination[0]][destination[1]] == math.MaxInt32 {
		return -1
	}
	return distance[destination[0]][destination[1]]
}

// shortestDistance_2 bfs
// 时间复杂度: O(m*n*max(m,n))
// 空间复杂度: O(m*n)
func shortestDistance_2(maze [][]int, start []int, destination []int) int {
	m, n := len(maze), len(maze[0])
	// distance[i][j]: 记录从起始位置到 (i, j) 的最小步数
	distance := make([][]int, m)
	for i := 0; i < m; i++ {
		distance[i] = make([]int, n)
		for j := 0; j < n; j++ {
			distance[i][j] = math.MaxInt32
		}
	}
	distance[start[0]][start[1]] = 0

	var (
		dirs  = [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
		queue [][]int
	)
	queue = append(queue, start)
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		for _, dir := range dirs {
			x, y := s[0]+dir[0], s[1]+dir[1]
			count := 0
			for x >= 0 && y >= 0 && x < m && y < n && maze[x][y] == 0 {
				x += dir[0]
				y += dir[1]
				count++
			}
			if distance[s[0]][s[1]]+count < distance[x-dir[0]][y-dir[1]] {
				distance[x-dir[0]][y-dir[1]] = distance[s[0]][s[1]] + count
				queue = append(queue, []int{x - dir[0], y - dir[1]})
			}
		}
	}
	if distance[destination[0]][destination[1]] == math.MaxInt32 {
		return -1
	}
	return distance[destination[0]][destination[1]]
}
