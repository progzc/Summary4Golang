package leetcode_1059_all_paths_from_source_lead_to_destination

// 1059. 从始点到终点的所有路径
// https://leetcode.cn/problems/all-paths-from-source-lead-to-destination/

// leadsToDestination dfs
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
// 缺点：最后一个用例会超时
func leadsToDestination(n int, edges [][]int, source int, destination int) bool {
	g := make(map[int][]int, 0)
	for _, edge := range edges {
		g[edge[0]] = append(g[edge[0]], edge[1])
	}
	// 若终点还有子节点, 直接返回false
	if _, ok := g[destination]; ok {
		return false
	}
	visit := make([]bool, n)
	visit[source] = true

	var dfs func(start, end int) bool
	dfs = func(start, end int) bool {
		if _, ok := g[start]; !ok {
			return start == end
		}
		for _, element := range g[start] {
			// 说明有环
			if visit[element] {
				return false
			}
			visit[element] = true
			if !dfs(element, end) {
				return false
			}
			visit[element] = false
		}
		return true
	}
	return dfs(source, destination)
}
