package leetcode_0121_best_time_to_buy_and_sell_stock

import "math"

// 0121.买卖股票的最佳时机🌟
// https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock/

// maxProfit 贪心法
// 时间复杂度：O(n)
// 空间复杂度：O(1)
func maxProfit(prices []int) int {
	minPrice, maxProfile := math.MaxInt64, 0
	for _, price := range prices {
		if price < minPrice {
			minPrice = price
			continue
		}
		if price-minPrice > maxProfile {
			maxProfile = price - minPrice
		}
	}
	return maxProfile
}

// maxProfit_2 动态规划
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
//			b.昨天不持股，今天买入股票（注意：只允许交易一次，因此手上的现金数就是当天的股价的相反数）。
func maxProfit_2(prices []int) int {
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
		// 注意: 写出下面这样会造成错误，这是因为dp[i-1][0]有两种情况:
		//	a.昨天不持股，今天什么都不做
		//	b.昨天持股，今天卖出股票（现金数增加）
		//dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
		dp[i][1] = max(dp[i-1][1], -prices[i])
	}
	return dp[n-1][0]
}

// maxProfit_3 动态规划（空间优化）
// 时间复杂度：O(n)
// 空间复杂度：O(1)
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
//			b.昨天不持股，今天买入股票（注意：只允许交易一次，因此手上的现金数就是当天的股价的相反数）。
func maxProfit_3(prices []int) int {
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
		dp[0] = max(dp[0], dp[1]+prices[i])
		// 注意: 写出下面这样会造成错误，这是因为dp[i-1][0]有两种情况:
		//	a.昨天不持股，今天什么都不做
		//	b.昨天持股，今天卖出股票（现金数增加）
		//dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
		dp[1] = max(dp[1], -prices[i])
	}
	return dp[0]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
