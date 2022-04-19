package leetcode_0377_combination_sum_iv

// 0377.组合总和Ⅳ
// https://leetcode-cn.com/problems/combination-sum-iv/

// combinationSum4 动态规划(组合背包)
// 时间复杂度: O(N*T)
// 空间复杂度: O(T)
// 思路：考虑顺序的组合背包
func combinationSum4(nums []int, target int) int {
	dp := make([]int, target+1)
	// 注意事项：只有当不选取任何元素时，元素之和才为0，因此只有1种方案
	dp[0] = 1
	for i := 1; i <= target; i++ {
		for j := 0; j < len(nums); j++ {
			if i >= nums[j] {
				dp[i] += dp[i-nums[j]]
			}
		}
	}
	return dp[target]
}
