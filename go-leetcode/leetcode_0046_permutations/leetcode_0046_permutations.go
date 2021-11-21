package leetcode_0046_permutations

// 46.全排列
// link: https://leetcode-cn.com/problems/permutations/

// permute 回溯
func permute(nums []int) [][]int {
	var ans [][]int
	n := len(nums)
	var dfs func(first int)
	dfs = func(first int) {
		// 所有的数都填完了
		if first == n {
			ans = append(ans, append([]int(nil), nums...))
		}
		for i := first; i < n; i++ {
			nums[first], nums[i] = nums[i], nums[first]
			dfs(first + 1)
			nums[first], nums[i] = nums[i], nums[first]
		}
	}
	dfs(0)
	return ans
}
