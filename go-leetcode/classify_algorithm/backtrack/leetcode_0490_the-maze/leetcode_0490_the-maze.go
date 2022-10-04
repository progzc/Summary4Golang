package leetcode_0490_the_maze

// 0490. 迷宫
// https://leetcode.cn/problems/the-maze/

// hasPath dfs
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*N)
func hasPath(maze [][]int, start []int, destination []int) bool {
	m, n := len(maze), len(maze[0])
	var (
		visited = make([][]bool, m)
		dfs     func(start, destination []int) bool
	)
	for i := 0; i < m; i++ {
		visited[i] = make([]bool, n)
	}

	dfs = func(start, destination []int) bool {
		if visited[start[0]][start[1]] {
			return false
		}
		if start[0] == destination[0] && start[1] == destination[1] {
			return true
		}
		visited[start[0]][start[1]] = true

		l, r, u, d := start[1]-1, start[1]+1, start[0]-1, start[0]+1

		// 向左
		for l >= 0 && maze[start[0]][l] == 0 {
			l--
		}
		if dfs([]int{start[0], l + 1}, destination) {
			return true
		}

		// 向右
		for r < n && maze[start[0]][r] == 0 {
			r++
		}
		if dfs([]int{start[0], r - 1}, destination) {
			return true
		}

		// 向上
		for u >= 0 && maze[u][start[1]] == 0 {
			u--
		}
		if dfs([]int{u + 1, start[1]}, destination) {
			return true
		}

		// 向下
		for d < m && maze[d][start[1]] == 0 {
			d++
		}
		if dfs([]int{d - 1, start[1]}, destination) {
			return true
		}

		return false
	}
	return dfs(start, destination)
}

// hasPath_2 bfs
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*N)
func hasPath_2(maze [][]int, start []int, destination []int) bool {
	m, n := len(maze), len(maze[0])
	var (
		visited = make([][]bool, m)
		dirs    = [][]int{{0, 1}, {0, -1}, {-1, 0}, {1, 0}}
		queue   [][]int
	)
	for i := 0; i < m; i++ {
		visited[i] = make([]bool, n)
	}
	queue = append(queue, start)
	visited[start[0]][start[1]] = true

	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		if s[0] == destination[0] && s[1] == destination[1] {
			return true
		}

		for _, dir := range dirs {
			x, y := s[0]+dir[0], s[1]+dir[1]
			for x >= 0 && y >= 0 && x < len(maze) && y < len(maze[0]) && maze[x][y] == 0 {
				x += dir[0]
				y += dir[1]
			}
			if !visited[x-dir[0]][y-dir[1]] {
				queue = append(queue, []int{x - dir[0], y - dir[1]})
				visited[x-dir[0]][y-dir[1]] = true
			}
		}
	}
	return false
}
