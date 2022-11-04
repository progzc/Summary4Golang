package leetcode_0096_unique_binary_search_trees

// 0096. 不同的二叉搜索树
// https://leetcode.cn/problems/unique-binary-search-trees/

// numTrees dfs（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func numTrees(n int) int {
	if n == 0 || n == 1 {
		return 1
	}

	count := 0
	for i := 1; i <= n; i++ {
		left := numTrees(i - 1)
		right := numTrees(n - i)
		count += left * right
	}
	return count
}

// numTrees_2 dfs+记忆法搜索（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func numTrees_2(n int) int {
	var (
		dfs func(n int) int
		m   = make(map[int]int, 0)
	)

	dfs = func(n int) int {
		if n == 0 || n == 1 {
			return 1
		}
		if v, ok := m[n]; ok {
			return v
		}

		count := 0
		for i := 1; i <= n; i++ {
			left := numTrees(i - 1)
			m[i-1] = left
			right := numTrees(n - i)
			m[n-i] = right
			count += left * right
		}
		m[n] = count
		return count
	}
	return dfs(n)
}

// numTrees_3 动态规划
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
// 思路:
//	状态: G[n]表示长度为 n 的序列能构成的不同二叉搜索树的个数。
//	初始状态: G[0]=1,G[1]=1
//	递推表达式: G[n] = Sum (G[i-1]*G[n-i]) i从1到n
func numTrees_3(n int) int {
	g := make([]int, n+1)
	g[0], g[1] = 1, 1
	for i := 2; i <= n; i++ {
		for j := 1; j <= i; j++ {
			g[i] += g[i-1] * g[i-j]
		}
	}
	return g[n]
}
