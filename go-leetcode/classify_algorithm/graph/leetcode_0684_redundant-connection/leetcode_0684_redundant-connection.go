package leetcode_0684_redundant_connection

// 0684. 冗余连接
// https://leetcode.cn/problems/redundant-connection/

// findRedundantConnection 并查集
func findRedundantConnection(edges [][]int) []int {
	n := len(edges)
	parent, rank := make([]int, n+1), make([]int, n+1)
	for i := 0; i < n+1; i++ {
		parent[i] = i
		rank[i] = 1
	}

	var (
		find  func(x int) int
		union func(x, y int) bool
	)
	find = func(x int) int {
		if parent[x] != x {
			// 路径压缩
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union = func(x, y int) (isRing bool) {
		fx, fy := find(x), find(y)
		if fx == fy {
			isRing = true
			return
		}

		// 按秩合并
		if rank[fx] <= rank[fy] {
			parent[fx] = fy
		} else {
			parent[fy] = fx
		}

		if rank[fx] == rank[fy] {
			rank[fx]++
		}
		return
	}

	for _, edge := range edges {
		if union(edge[0], edge[1]) {
			return edge
		}
	}

	return nil
}
