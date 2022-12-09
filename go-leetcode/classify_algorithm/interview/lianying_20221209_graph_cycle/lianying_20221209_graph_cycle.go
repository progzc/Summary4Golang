package lianying_20221209_graph_cycle

func hasCycle(graph [][]int) bool {
	m := make(map[int][]int)
	rows, columns := len(graph), len(graph[0])
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if graph[i][j] == 1 {
				m[i] = append(m[i], j)
			}
		}
	}

	visited := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		visited[i] = make([]bool, columns)
	}

	var dfs func(i, j int) bool
	dfs = func(i, j int) bool {
		if graph[i][j] == 0 {
			return false
		}
		if visited[i][j] {
			return true
		}
		visited[i][j] = true
		defer func() {
			visited[i][j] = false
		}()
		for _, next := range m[j] {
			if dfs(j, next) {
				return true
			}
		}
		return false
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if graph[i][j] == 1 {
				if dfs(i, j) {
					return true
				}
			}
		}
	}
	return false
}
