package leetcode_0787_cheapest_flights_within_k_stops

// 0787. K 站中转内最便宜的航班
// https://leetcode.cn/problems/cheapest-flights-within-k-stops/

// findCheapestPrice Bellman-Ford
// 时间复杂度: O(k*(n+m))
// 空间复杂度: O(n)
// Bellman-Ford算法的讲解：https://zhuanlan.zhihu.com/p/352724346
func findCheapestPrice(n int, flights [][]int, src int, dst int, k int) int {
	const INF = 0x3f3f3f3f
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = INF
	}
	dist[src] = 0
	for limit := 0; limit < k+1; limit++ {
		clone := make([]int, len(dist))
		copy(clone, dist)
		for _, flight := range flights {
			x, y, w := flight[0], flight[1], flight[2]
			dist[y] = min(dist[y], clone[x]+w)
		}
	}
	// 这里不像Dijkstra写等于正无穷是因为可能有负权边甚至是负环的存在，使得"正无穷"在迭代过程中受到一点影响
	if dist[dst] > INF/2 {
		return -1
	}
	return dist[dst]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
