package leetcode_0518_coin_change_2

// 0518.零钱兑换 II
// https://leetcode-cn.com/problems/coin-change-2/

// change 动态规划(完全背包的组合问题)
// 时间复杂度：O(A*C)
// 空间复杂度：O(A)
// 特点：完全背包不考虑顺序的组合问
func change(amount int, coins []int) int {
	// dp[i]表示金额之和等于i的硬币组合数
	dp := make([]int, amount+1)
	// 只有当不选取任何硬币时，金额之和才为0，因此只有1种硬币组合
	dp[0] = 1
	for _, coin := range coins {
		for i := 1; i <= amount; i++ {
			if i >= coin {
				dp[i] = dp[i] + dp[i-coin]
			}
		}
	}
	return dp[amount]
}
