package leetcode_0122_best_time_to_buy_and_sell_stock_ii

import "math"

// 0122.买卖股票的最佳时机 II🌟
// https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock-ii/

// maxProfit 贪心法
// 时间复杂度：O(n)
// 空间复杂度: O(1)
func maxProfit(prices []int) int {
	minPrice, maxProfile := math.MaxInt64, 0
	for _, price := range prices {
		if price < minPrice {
			minPrice = price
		} else if price > minPrice {
			maxProfile += price - minPrice
			minPrice = price
		}
	}
	return maxProfile
}

// maxProfit_4 动态规划
// 时间复杂度：O(n)
// 空间复杂度：O(n)
// 思路:
//
//	状态: dp[i][j]表示天数 [0,i] 区间里，下标 i 这一天状态为 j 的时候持有的现金数。
//		 其中, j = 0, 表示当前不持股; j = 1, 表示当前持股。
//	转移方程:
//		dp[i][0]: 规定了今天不持股，有以下两种情况:
//			a.昨天不持股，今天什么都不做
//			b.昨天持股，今天卖出股票（现金数增加）
//		dp[i][1]: 规定了今天持股，有以下两种情况:
//			a.昨天持股，今天什么都不做（现金数与昨天一样）。
//			b.昨天不持股，今天买入股票（注意：这里与【0121.买卖股票的最佳时机】有区别）。
func maxProfit_4(prices []int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, 2)
	}

	// dp[i][0] 下标为 i 这天结束的时候，不持股，手上拥有的现金数
	// dp[i][1] 下标为 i 这天结束的时候，持股，手上拥有的现金数
	// 初始化：不持股显然为 0，持股就需要减去第 1 天（下标为 0）的股价
	dp[0][0] = 0
	dp[0][1] = -prices[0]

	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i])
		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
	}
	return dp[n-1][0]
}

// maxProfit_5 动态规划（空间优化）
// 时间复杂度：O(n)
// 空间复杂度：O(n)
// 思路:
//
//	状态: dp[i][j]表示天数 [0,i] 区间里，下标 i 这一天状态为 j 的时候持有的现金数。
//		 其中, j = 0, 表示当前不持股; j = 1, 表示当前持股。
//	转移方程:
//		dp[i][0]: 规定了今天不持股，有以下两种情况:
//			a.昨天不持股，今天什么都不做
//			b.昨天持股，今天卖出股票（现金数增加）
//		dp[i][1]: 规定了今天持股，有以下两种情况:
//			a.昨天持股，今天什么都不做（现金数与昨天一样）。
//			b.昨天不持股，今天买入股票（注意：这里与【0121.买卖股票的最佳时机】有区别）。
func maxProfit_5(prices []int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}
	dp := make([]int, 2)

	// dp[i][0] 下标为 i 这天结束的时候，不持股，手上拥有的现金数
	// dp[i][1] 下标为 i 这天结束的时候，持股，手上拥有的现金数

	// 初始化：不持股显然为 0，持股就需要减去第 1 天（下标为 0）的股价
	dp[0] = 0
	dp[1] = -prices[0]

	for i := 1; i < n; i++ {
		dp0, dp1 := dp[0], dp[1]
		dp[0] = max(dp0, dp1+prices[i])
		dp[1] = max(dp1, dp0-prices[i])
	}
	return dp[0]
}

// maxProfit_2 贪心法(优化)
// 时间复杂度：O(n)
// 空间复杂度: O(1)
func maxProfit_2(prices []int) int {
	maxProfile := 0
	for i := 1; i < len(prices); i++ {
		maxProfile += max(0, prices[i]-prices[i-1])
	}
	return maxProfile
}

// maxProfit_3 dfs(超时)
// 时间复杂度：O(n^2)
// 空间复杂度: O(n)
func maxProfit_3(prices []int) int {
	var (
		ans int
		dfs func(idx, pre, sum int, flag bool)
	)

	n := len(prices)
	dfs = func(idx, pre, sum int, flag bool) {
		if idx == n {
			return
		}
		if !flag {
			// 未持有股票
			// 选择买入
			dfs(idx+1, prices[idx], sum, true)
			// 选择不买入
			dfs(idx+1, pre, sum, false)
		} else {
			// 持有股票
			// 选择卖出
			ans = max(ans, sum+prices[idx]-pre)
			dfs(idx+1, 0, sum+prices[idx]-pre, false)
			// 选择不卖出
			dfs(idx+1, pre, sum, true)
		}
	}
	dfs(0, 0, 0, false)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
