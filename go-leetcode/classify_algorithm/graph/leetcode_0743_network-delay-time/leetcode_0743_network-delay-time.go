package leetcode_0743_network_delay_time

import "math"

// 0743. 网络延迟时间
// https://leetcode.cn/problems/network-delay-time/

// 几篇比较好的总结：
//	吃透Dijkstra：https://leetcode.cn/problems/network-delay-time/solution/gtalgorithm-dan-yuan-zui-duan-lu-chi-tou-w3zc/
//	一题七解：https://leetcode.cn/problems/network-delay-time/solution/by-tong-zhu-ud0k/
//	五种最短路径算法总结：https://leetcode.cn/problems/network-delay-time/solution/dirkdtra-by-happysnaker-vjii/

// networkDelayTime Dijkstra最短路径算法
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
// 思路：
//	根据题意，从节点 k 发出的信号，到达节点 x 的时间就是节点 k 到节点 x 的最短路的长度。
//	因此我们需要求出节点 k 到其余所有点的最短路，其中的最大值就是答案。若存在从 k 出发无法到达的点，则返回 −1。
// 注意事项：
//	下面的代码将节点编号减小了 1，从而使节点编号位于 [0,n-1] 范围。
func networkDelayTime(times [][]int, n int, k int) int {
	// 邻接矩阵
	const inf = math.MaxInt32 >> 2
	g := make([][]int, n)
	for i := range g {
		g[i] = make([]int, n)
		for j := range g[i] {
			g[i][j] = inf
		}
	}
	for _, t := range times {
		x, y := t[0]-1, t[1]-1
		g[x][y] = t[2]
	}

	// Dijkstra 初始化
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[k-1] = 0
	used := make([]bool, n)
	// Dijkstra 搜索
	for i := 0; i < n; i++ {
		// 从 未确定节点 中取一个与起点节点距离最近的点，并将其归类为 已确定节点
		x := -1
		for y, u := range used {
			if !u && (x == -1 || dist[y] < dist[x]) {
				x = y
			}
		}
		used[x] = true

		// 根据 已确定节点 去更新其他节点
		for y, time := range g[x] {
			dist[y] = min(dist[y], dist[x]+time)
		}
	}

	ans := 0
	for _, d := range dist {
		if d == inf {
			return -1
		}
		ans = max(ans, d)
	}
	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
