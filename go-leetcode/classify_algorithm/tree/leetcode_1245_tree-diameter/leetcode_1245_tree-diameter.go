package leetcode_1245_tree_diameter

// 1245. 树的直径
// https://leetcode.cn/problems/tree-diameter/

// treeDiameter 树形dfs
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：以任意一点为root，与其关联的点为子节点，求出所有子节点对应的路径的前两大值max1,max2，则该点对应的最长路径为max1+max2
func treeDiameter(edges [][]int) int {
	g := make(map[int][]int)
	for _, edge := range edges {
		g[edge[0]] = append(g[edge[0]], edge[1])
		g[edge[1]] = append(g[edge[1]], edge[0])
	}

	var (
		visited = make([]bool, len(edges)+1)
		dfs     func(idx int) int
		ans     int
	)
	dfs = func(idx int) int {
		visited[idx] = true
		max1, max2 := 0, 0
		for _, next := range g[idx] {
			if !visited[next] {
				num := dfs(next)
				if num > max1 {
					max2 = max1
					max1 = num
				} else if num > max2 {
					max2 = num
				}
			}
		}
		ans = max(ans, max1+max2)
		return max(max1, max2) + 1
	}
	dfs(0)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
