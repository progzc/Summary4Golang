package leetcode_0343_integer_break

// 343. 整数拆分
// https://leetcode.cn/problems/integer-break/

// integerBreak 动态规划
// 时间复杂度: O()
// 空间复杂度: O()
// 思路:
//	状态: dp[i]表示使正整数i拆分后乘积的最大值
//	转移方程:
//		当n>=2时, dp[i] = max{j*(i-j),j*dp[i-j]},其中j属于[1...i-1]
//	含义: 当 i ≥2 时，假设对正整数 i 拆分出的第一个正整数是 j (1≤j<i)，则有以下两种方案：
//		将 i 拆分成 j 和 i-j 的和，且 i-j 不再拆分成多个正整数，此时的乘积是 j×(i−j)
//		将 i 拆分成 j 和 i-j 的和，且 i-j 继续拆分成多个正整数，此时的乘积是 j×dp[i−j]
//	初始值:
//		dp[0] = 0
//		dp[1] = 0
func integerBreak(n int) int {
	dp := make([]int, n+1)
	for i := 2; i < n+1; i++ {
		for j := 1; j < i; j++ {
			dp[i] = max(dp[i], max(j*dp[i-j], j*(i-j)))
		}
	}
	return dp[n]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
