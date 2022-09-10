package leetcode_0839_similar_string_groups

// 0839. 相似字符串组
// https://leetcode.cn/problems/similar-string-groups/

// numSimilarGroups 并查集
// 时间复杂度: O(n^2*m+n*log(n))
// 空间复杂度: O(n)
func numSimilarGroups(strs []string) int {
	n := len(strs)
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	var (
		find      func(x int) int
		union     func(x, y int)
		isSimilar func(s, t string) bool
	)
	find = func(x int) int {
		if parent[x] != x {
			// 压缩路径
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union = func(x, y int) {
		fx, fy := find(x), find(y)
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

	// 判断两个字符串相似
	// 注意前提条件：s和t具有相同的长度,且是彼此的字母异位词
	isSimilar = func(s, t string) bool {
		diff := 0
		for i := range s {
			if s[i] != t[i] {
				diff++
				if diff > 2 {
					return false
				}
			}
		}
		return true
	}

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if find(i) != find(j) && isSimilar(strs[i], strs[j]) {
				union(i, j)
			}
		}
	}

	count := 0
	for i := range parent {
		if parent[i] == i {
			count++
		}
	}

	return count
}
