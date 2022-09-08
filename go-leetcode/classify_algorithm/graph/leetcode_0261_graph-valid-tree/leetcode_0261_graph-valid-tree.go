package leetcode_0261_graph_valid_tree

// 0261. 以图判树
// https://leetcode.cn/problems/graph-valid-tree/

// 图的一些概念和性质：https://blog.csdn.net/zsy3757486/article/details/125266607
// 判断图是否有环的三种方法：https://leetcode.cn/problems/graph-valid-tree/solution/java-pan-duan-tu-zhong-you-huan-de-san-chong-fang-/
//	a.并查集
//	b.BFS
//	c.DFS

// validTree 并查集
// 思路：树需要满足的两个条件是：没有环；连通分量为1
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func validTree(n int, edges [][]int) bool {
	// p 用来存储每个节点的祖宗
	p := make([]int, n)
	// 初始化，每个节点的祖宗初始化为自己
	for i := 0; i < n; i++ {
		p[i] = i
	}

	// 集合的合并操作
	for i := 0; i < len(edges); i++ {
		// 如果出现可环,则直接退出
		if find(edges[i][0], p) == find(edges[i][1], p) {
			return false
		} else {
			union(edges[i][0], edges[i][1], p)
		}
	}

	// 查询连通分量
	c := 0
	for i := 0; i < n; i++ {
		// 根节点的个数其实就等于连通分量
		if p[i] == i {
			c++
		}
	}

	// 判断连通分量是否为1（即判断是否只有一个根节点）
	return c == 1
}

// union 两个集合的合并操作
func union(x, y int, p []int) {
	// 若两个节点的祖宗不相等,则将两个集合合并
	if fx, fy := find(x, p), find(y, p); fx != fy {
		p[fx] = fy
	}
}

// find 查找一个节点的祖宗节点
func find(x int, p []int) int {
	// 如果当前节点的祖宗为其本身,则为根节点
	if p[x] == x {
		return x
	}
	// 如果不是，则递归查找
	return find(p[x], p)
}
