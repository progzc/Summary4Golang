package leetcode_0547_number_of_provinces

// 0547. 省份数量
// https://leetcode.cn/problems/number-of-provinces/

// findCircleNum 并查集（按秩合并+路径压缩）
// 时间复杂度: O(n^2*log(n))
// 空间复杂度: O(n)
func findCircleNum(isConnected [][]int) int {
	n := len(isConnected)
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	var (
		find  func(x int) int
		union func(x, y int)
	)
	// 查找
	find = func(x int) int {
		if parent[x] != x {
			// 路径压缩
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	// 按秩合并
	union = func(x, y int) {
		if fx, fy := find(x), find(y); fx != fy {
			if rank[fx] <= rank[fy] {
				parent[fx] = fy
			} else {
				parent[fy] = fx
			}

			if rank[fx] == rank[fy] {
				rank[fx]++
			}
		}
	}

	for i, row := range isConnected {
		for j := i + 1; j < n; j++ {
			if row[j] == 1 {
				union(i, j)
			}
		}
	}

	// 计算连通分量
	count := 0
	for i := range parent {
		if parent[i] == i {
			count++
		}
	}
	return count
}
