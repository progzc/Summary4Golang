package leetcode_0518_coin_change_2

// 0518.零钱兑换 II
// https://leetcode-cn.com/problems/coin-change-2/

// change 动态规划
// 时间复杂度：O()
// 空间复杂度：O()
// 思路：动态规划
//	定义：dp[i]表示金额之和等于i的硬币组合数
func change(amount int, coins []int) int {
	dp := make([]int, amount+1)
	// 只有当不选取任何硬币时，金额之和才为0，因此只有1种硬币组合
	dp[0] = 1
	for _, coin := range coins {
		for i := coin; i <= amount; i++ {
			dp[i] = dp[i] + dp[i-coin]
		}
	}
	return dp[amount]
}
