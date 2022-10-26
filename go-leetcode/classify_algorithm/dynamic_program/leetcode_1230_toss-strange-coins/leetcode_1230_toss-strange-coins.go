package leetcode_1230_toss_strange_coins

// 1230. 抛掷硬币
// https://leetcode.cn/problems/toss-strange-coins/

// probabilityOfHeads 动态规划
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路：
// 	状态定义: dp[i][j]表示从0到第i枚硬币(含第i枚),正面朝上个数等于j的概率
//	转移放出: dp[i][j] = dp[i-1][j-1]*prob[i] + dp[i-1][j]*(1-prob[i])
func probabilityOfHeads(prob []float64, target int) float64 {
	n := len(prob)
	dp := make([][]float64, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]float64, target+1)
		if i > 0 {
			dp[i][0] = dp[i-1][0] * (1 - prob[i])
		} else {
			dp[0][0] = 1 - prob[0]
			if target > 0 {
				dp[0][1] = prob[0]
			}
		}
	}

	for i := 1; i < n; i++ {
		for j := 1; j < target+1; j++ {
			dp[i][j] = dp[i-1][j-1]*prob[i] + dp[i-1][j]*(1-prob[i])
		}
	}
	return dp[n-1][target]
}
