package leetcode_0879_profitable_schemes

// 0879.盈利计划
// https://leetcode-cn.com/problems/profitable-schemes/

// profitableSchemes 动态规划（多维背包）
// 时间复杂度：O(G*N*M)
// 空间复杂度：O(N*M)
// 思路：
//	a.三维背包问题：
//		状态：dp[i][j][k]：表示在前i个工作中选择了j个员工，并且满足工作利润至少为k的情况下的盈利计划的总数目
//		初始值：
//			i) 初始化为 dp[0][0][0] = 1，dp数组表示的意思是，进行前 i 种工作，恰好使用 j 个人，工作利润至少为 k 的情况数量。
//				注意：这种最后需要对结果进行累加
//			ii) 初始化整个 dp[0][j][0] = 1，dp数组表示的意思是，进行前 i 种工作，使用 j 个人，工作利润至少为 k 的情况数量。（推荐这种初始化条件）
//	b.二维背包问题：
//		状态：dp[j][k]：表示在前i个工作中选择了j个员工，并且满足工作利润至少为k的情况下的盈利计划的总数目
//		初始值：dp[j][0]=1, 表示进行前 i 种工作，使用 j 个人，工作利润至少为 k 的情况数量。
func profitableSchemes(n int, minProfit int, group []int, profit []int) int {
	const mod = 1e9 + 7
	// 初始化
	dp := make([][]int, n+1)
	for j := range dp {
		dp[j] = make([]int, minProfit+1)
		dp[j][0] = 1
	}
	for i, member := range group {
		earn := profit[i]
		// 注意下面都要逆序，这样才可以使用滚动数组
		for j := n; j >= member; j-- {
			for k := minProfit; k >= 0; k-- {
				// 注意是利润至少为k，对于dp[j-member][max(0, k-earn)]的解释如下：
				//		若earn > k, 则当前利润至少为earn了，剩下还需要利润至少为0
				//		若earn < k, 则当前利润为earn， 剩下还需要利润至少为 k-earn
				// 这种写法是错误的：dp[j][k] = (dp[j][k] + dp[j-member][k-earn] % mod
				dp[j][k] = (dp[j][k] + dp[j-member][max(0, k-earn)]) % mod
			}
		}
	}
	return dp[n][minProfit]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
