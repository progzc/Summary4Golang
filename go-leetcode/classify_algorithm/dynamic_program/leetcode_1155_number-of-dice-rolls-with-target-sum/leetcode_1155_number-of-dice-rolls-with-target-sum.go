package leetcode_1155_number_of_dice_rolls_with_target_sum

// 1155.掷骰子的N种方法
// https://leetcode-cn.com/problems/number-of-dice-rolls-with-target-sum/

// numRollsToTarget 动态规划(分组0/1背包)
// 时间复杂度: O(N*K*T)
// 空间复杂度: O(N*T)
// 思路：转化为分组0/1背包的组合问题
func numRollsToTarget(n int, k int, target int) int {
	// 注意事项：防止溢出
	// 答案可能很大，你需要对 10^9 + 7 取模
	// 为什么需要对1000000007取模？因为1000000007足够大，且为质数, 且能保证加法和乘法操作时不会溢出
	mod := 1000000007

	//dp[i][j]表示投掷i个骰子点数和为j的方法数
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, target+1)
	}
	//投掷0个骰子点数和为0的方法数为1
	//投掷超过1个筛子点数和为0的方法数都为0
	dp[0][0] = 1

	for i := 1; i <= n; i++ {
		for j := 1; j <= target; j++ {
			for m := 1; m <= k && j >= m; m++ {
				dp[i][j] = (dp[i][j] + dp[i-1][j-m]) % mod
			}
		}
	}
	return dp[n][target]
}
