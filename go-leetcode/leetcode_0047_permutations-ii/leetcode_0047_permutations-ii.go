package leetcode_0047_permutations_ii

import "sort"

// 47.全排列II
// link: https://leetcode-cn.com/problems/permutations-ii/

// permuteUnique 回溯法
// 时间复杂度 O(n*n!)
// 空间复杂度 O(n)
func permuteUnique(nums []int) [][]int {
	sort.Ints(nums)

	n := len(nums)
	used := make([]bool, n)
	var ans [][]int
	var perm []int
	var dfs func(depth int)
	dfs = func(depth int) {
		if depth == n {
			ans = append(ans, append([]int(nil), perm...))
			return
		}
		for i := 0; i < n; i++ {
			if used[i] {
				continue
			}
			// 剪枝
			if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
				continue
			}
			perm = append(perm, nums[i])
			used[i] = true
			dfs(depth + 1)
			used[i] = false
			perm = perm[:len(perm)-1]
		}
	}
	dfs(0)
	return ans
}
