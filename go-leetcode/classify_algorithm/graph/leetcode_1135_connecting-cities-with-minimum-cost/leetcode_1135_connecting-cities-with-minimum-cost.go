package leetcode_1135_connecting_cities_with_minimum_cost

import (
	"sort"
)

// 1135. 最低成本联通所有城市
// https://leetcode.cn/problems/connecting-cities-with-minimum-cost/

// minimumCost 排序+并查集
// 时间复杂度: O(m*log(m)+m*n)
// 空间复杂度: O(n)
func minimumCost(n int, connections [][]int) int {
	// 初始化
	parent := make([]int, n+1)
	rank := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		parent[i] = i
		rank[i] = 1
	}

	// 排序
	sort.Slice(connections, func(i, j int) bool {
		return connections[i][2] < connections[j][2]
	})

	count, cost := 0, 0
	for _, conn := range connections {
		// 已经有n-1条边,说明所有的点已经联通
		if count == n-1 {
			break
		}
		// 如果成环,则不需要加入
		if find(conn[0], parent) == find(conn[1], parent) {
			continue
		}
		// 不成环,则合并
		union(conn[0], conn[1], parent, rank)
		count++
		cost += conn[2]
	}

	if count != n-1 {
		return -1
	}
	return cost
}

func union(x, y int, parent, rank []int) {
	fx, fy := find(x, parent), find(y, parent)
	if fx != fy {
		// 按秩合并
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

func find(x int, parent []int) int {
	if parent[x] != x {
		// 路径压缩
		parent[x] = find(parent[x], parent)
	}
	return parent[x]
}
