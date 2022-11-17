package leetcode_0377_combination_sum_iv

// 0377.组合总和Ⅳ
// https://leetcode-cn.com/problems/combination-sum-iv/

// combinationSum4 动态规划(组合背包)一维动态规划
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

// combinationSum4 动态规划(组合背包)二维动态规划
// 时间复杂度: O(N*T^2)
// 空间复杂度: O(T)
// 思路：考虑顺序的组合背包
// 状态: dp[i][j]表示组合长度为i,凑成总和为j的方案数。由于对组合方案的长度没有限制，因此我们最终答案为所有的 f[x][target] 的总和。
// 初始值: dp[0][0] = 1
// 转移方程: dp[len][target] 是以下所有方案的总和，前置条件是target>=nums[i]
//		1.最后一个数选择nums[0]，方案数为dp[len-1][target-nums[0]]
//		2.最后一个数选择nums[1]，方案数为dp[len-1][target-nums[1]]
//		3.最后一个数选择nums[2]，方案数为dp[len-1][target-nums[2]]
//		....
func combinationSum4_2(nums []int, target int) int {
	// 因为 nums[i] 最小值为 1，因此构成答案的最大长度为 target
	l := target
	dp := make([][]int, l+1)
	for i := 0; i < l+1; i++ {
		dp[i] = make([]int, target+1)
	}

	// 初始化
	dp[0][0] = 1
	ans := 0
	for i := 1; i < l+1; i++ {
		for j := 0; j < target+1; j++ {
			for _, num := range nums {
				if num <= j {
					dp[i][j] += dp[i-1][j-num]
				}
			}
		}
		ans += dp[i][target]
	}
	return ans
}
