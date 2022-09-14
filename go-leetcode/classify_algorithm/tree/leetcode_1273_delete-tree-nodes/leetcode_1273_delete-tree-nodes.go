package leetcode_1273_delete_tree_nodes

// 1273. 删除树节点
// https://leetcode.cn/problems/delete-tree-nodes/

// deleteTreeNodes 构造树+dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func deleteTreeNodes(nodes int, parent []int, value []int) int {
	// 构造树
	g := make(map[int][]int)
	for i, p := range parent {
		if p != -1 {
			g[p] = append(g[p], i)
		}
	}

	counter := make([]int, nodes)
	for i := range counter {
		counter[i] = 1
	}

	// dfs
	var dfs func(u int)
	dfs = func(u int) {
		if vs, ok := g[u]; ok {
			for _, v := range vs {
				dfs(v)
				value[u] += value[v]
				counter[u] += counter[v]
			}
		}

		if value[u] == 0 {
			counter[u] = 0
		}
	}

	dfs(0)
	return counter[0]
}
