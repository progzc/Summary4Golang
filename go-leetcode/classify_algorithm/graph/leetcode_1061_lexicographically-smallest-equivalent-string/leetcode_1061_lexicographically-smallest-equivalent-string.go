package leetcode_1061_lexicographically_smallest_equivalent_string

// 1061. 按字典序排列最小的等效字符串
// https://leetcode.cn/problems/lexicographically-smallest-equivalent-string/

// smallestEquivalentString 并查集
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(n)
func smallestEquivalentString(s1 string, s2 string, baseStr string) string {
	parent := make([]int, 26)
	for i := 0; i < 26; i++ {
		parent[i] = i
	}

	var (
		find  func(x int) int
		union func(x, y int)
	)

	find = func(x int) int {
		if parent[x] != x {
			// 路径压缩
			parent[x] = find(parent[x])
		}
		return parent[x]
	}

	union = func(x, y int) {
		fx, fy := find(x), find(y)
		if fx != fy {
			// 按值大小合并
			if fx < fy {
				parent[fy] = fx
			} else {
				parent[fx] = fy
			}
		}
	}

	n, b := len(s1), len(baseStr)
	for i := 0; i < n; i++ {
		union(int(s1[i]-'a'), int(s2[i]-'a'))
	}

	ans := make([]byte, b)
	for i := 0; i < b; i++ {
		root := find(int(baseStr[i] - 'a'))
		ans[i] = byte(root + 'a')
	}
	return string(ans)
}
