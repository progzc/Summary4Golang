package leetcode_0309_best_time_to_buy_and_sell_stock_with_cooldown

// 309. 最佳买卖股票时机含冷冻期
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-with-cooldown/

// maxProfit dfs（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func maxProfit(prices []int) int {
	var (
		ans int
		dfs func(idx, pre, sum int, hold, preSoldOut bool)
	)
	n := len(prices)
	dfs = func(idx, pre, sum int, hold, preSoldOut bool) {
		if idx == n {
			return
		}
		if !hold {
			// 如果不持有股票
			// 选择买入
			if !preSoldOut {
				dfs(idx+1, prices[idx], sum, true, false)
			}
			// 选择不买入
			dfs(idx+1, pre, sum, false, false)
		} else {
			// 如果持有股票
			// 选择卖出
			if prices[idx] > pre {
				ans = max(ans, sum+prices[idx]-pre)
				dfs(idx+1, 0, sum+prices[idx]-pre, false, true)
			}
			// 选择不卖出
			dfs(idx+1, pre, sum, true, false)
		}
	}
	dfs(0, 0, 0, false, false)
	return ans
}

// maxProfit_2
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	状态: dp[i][j]表示第i天结束之后的持有的现金数，其中j的含义如下：
//		 j = 0: 表示目前持有一只股票
//		 j = 1: 表示目前不持有一只股票,且处于冷冻期。这里的「处于冷冻期」指的是在第 i 天结束之后的状态。
//	       	    也就是说：如果第 i 天结束之后处于冷冻期，那么第 i+1 天无法买入股票。
//		 j = 2: 表示目前不持有一只股票,且不处于冷冻期
//	转移方程:
//		dp[i][0] = max(dp[i-1][0], dp[i-1][2]-prices[i])
//		dp[i][1] = dp[i-1][0]+prices[i]
//		dp[i][2] = max(dp[i-1][1], dp[i-1][2])
//	边界条件:
//		dp[0][0] = -prices[0]
//		dp[0][1] = 0
//		dp[0][2] = 0
func maxProfit_2(prices []int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}

	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, 3)
	}
	// 初始边界条件
	dp[0][0] = -prices[0]
	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][0], dp[i-1][2]-prices[i])
		dp[i][1] = dp[i-1][0] + prices[i]
		dp[i][2] = max(dp[i-1][1], dp[i-1][2])
	}
	return max(max(dp[n-1][0], dp[n-1][1]), dp[n-1][2])
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
