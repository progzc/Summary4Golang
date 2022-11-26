package leetcode_0714_best_time_to_buy_and_sell_stock_with_transaction_fee

// 714. 买卖股票的最佳时机含手续费
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-with-transaction-fee/

// maxProfit dfs（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func maxProfit(prices []int, fee int) int {
	var (
		ans int
		dfs func(idx, pre, sum int, hold bool)
	)

	n := len(prices)
	if n < 2 {
		return ans
	}
	dfs = func(idx, pre, sum int, hold bool) {
		if idx == n {
			return
		}
		if !hold {
			// 当前未持有股票
			// 选择买入
			dfs(idx+1, prices[idx], sum, true)
			// 选择不买入
			dfs(idx+1, pre, sum, false)
		} else {
			// 当前持有股票
			// 选择卖出
			if prices[idx] > pre+fee {
				ans = max(ans, sum+prices[idx]-pre-fee)
				dfs(idx+1, 0, sum+prices[idx]-pre-fee, false)
			}
			// 选择不卖出
			dfs(idx+1, pre, sum, true)
		}
	}
	dfs(0, 0, 0, false)
	return ans
}

// maxProfit_2 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	状态: dp[i][j]表示第i天结束之后的持有的现金数, 其中j的含义如下:
//		 j = 0: 表示目前未持有一只股票
//		 j = 1: 表示目前持有一只股票
//	转移方程:
//		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i]-fee)
//		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
//	边界条件:
//		dp[0][0] = 0
//		dp[0][1] = -prices[0]
func maxProfit_2(prices []int, fee int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}

	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, 2)
	}
	// 边界条件
	dp[0][1] = -prices[0]

	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i]-fee)
		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
	}
	return dp[n-1][0]
}

// maxProfit_3 动态规划（空间优化）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路:
//	状态: dp[i][j]表示第i天结束之后的持有的现金数, 其中j的含义如下:
//		 j = 0: 表示目前未持有一只股票
//		 j = 1: 表示目前持有一只股票
//	转移方程:
//		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i]-fee)
//		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
//	边界条件:
//		dp[0][0] = 0
//		dp[0][1] = -prices[0]
func maxProfit_3(prices []int, fee int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}

	dp := make([]int, 2)
	// 边界条件
	dp[1] = -prices[0]

	for i := 1; i < n; i++ {
		dp0, dp1 := dp[0], dp[1]
		dp[0] = max(dp0, dp1+prices[i]-fee)
		dp[1] = max(dp1, dp0-prices[i])
	}
	return dp[0]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
