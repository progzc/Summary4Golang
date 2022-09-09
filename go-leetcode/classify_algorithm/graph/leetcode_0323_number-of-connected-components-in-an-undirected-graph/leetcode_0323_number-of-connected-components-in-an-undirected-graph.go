package leetcode_0323_number_of_connected_components_in_an_undirected_graph

// 0323. 无向图中连通分量的数目
// https://leetcode.cn/problems/number-of-connected-components-in-an-undirected-graph/

// countComponents 并查集
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func countComponents(n int, edges [][]int) int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}

	for i := 0; i < len(edges); i++ {
		if find(edges[i][0], p) != find(edges[i][1], p) {
			union(edges[i][0], edges[i][1], p)
		}
	}

	c := 0
	for i := 0; i < n; i++ {
		if p[i] == i {
			c++
		}
	}
	return c
}

func union(x, y int, p []int) {
	if fx, fy := find(x, p), find(y, p); fx != fy {
		p[fx] = fy
	}
}

func find(x int, p []int) int {
	if p[x] == x {
		return x
	}
	return find(p[x], p)
}
