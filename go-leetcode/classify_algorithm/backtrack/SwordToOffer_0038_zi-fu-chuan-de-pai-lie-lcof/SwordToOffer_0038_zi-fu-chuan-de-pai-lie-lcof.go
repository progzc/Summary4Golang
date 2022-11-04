package SwordToOffer_0038_zi_fu_chuan_de_pai_lie_lcof

import "sort"

// 剑指 Offer 38. 字符串的排列
// https://leetcode.cn/problems/zi-fu-chuan-de-pai-lie-lcof/

// 相关题目:
// 0047.全排列 II: https://leetcode-cn.com/problems/permutations-ii/

// permutation dfs
// 时间复杂度：O(n*n!)
// 空间复杂度：O(2n)
func permutation(s string) []string {
	var ans []string
	if len(s) == 0 {
		return ans
	}

	t := []byte(s)
	n := len(t)
	sort.Slice(t, func(i, j int) bool {
		return t[i] < t[j]
	})

	var (
		dfs  func(idx int, output []byte)
		used = make([]bool, n)
	)
	dfs = func(idx int, output []byte) {
		if idx == n {
			ans = append(ans, string(output))
			return
		}

		for i := 0; i < n; i++ {
			if used[i] || (i > 0 && !used[i-1] && t[i] == t[i-1]) {
				continue
			}

			output = append(output, t[i])
			used[i] = true
			dfs(idx+1, output)
			used[i] = false
			output = output[:len(output)-1]
		}
	}
	dfs(0, nil)
	return ans
}
