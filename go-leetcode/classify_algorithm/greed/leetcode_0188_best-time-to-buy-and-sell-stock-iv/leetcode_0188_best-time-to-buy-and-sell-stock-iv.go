package leetcode_0188_best_time_to_buy_and_sell_stock_iv

import "math"

// 188. 买卖股票的最佳时机 IV
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iv/

// maxProfit dfs（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func maxProfit(k int, prices []int) int {
	var (
		ans int
		dfs func(idx, count, pre, sum int, flag bool)
	)

	n := len(prices)
	dfs = func(idx, count, pre, sum int, flag bool) {
		if idx == n || count == k {
			return
		}
		if !flag {
			// 当前不持有股票
			// 选择买入当前股票
			dfs(idx+1, count, prices[idx], sum, true)
			// 选择不买入当前股票
			dfs(idx+1, count, pre, sum, false)
		} else {
			// 当前持有股票
			// 选择卖出股票
			ans = max(ans, sum+prices[idx]-pre)
			dfs(idx+1, count+1, 0, sum+prices[idx]-pre, false)
			// 选择不卖出股票
			dfs(idx+1, count, pre, sum, true)
		}
	}
	dfs(0, 0, 0, 0, false)
	return ans
}

// maxProfit_2 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	状态:
//		buy[i][j] 表示对于数组 prices[0..i] 中的价格而言，进行恰好 j 笔交易，并且当前手上持有一支股票，这种情况下的最大利润。
//		sell[i][j] 表示恰好进行 j 笔交易，并且当前手上不持有股票，这种情况下的最大利润。
//	转移方程:
//		buy[i][j] = max{buy[i−1][j],sell[i−1][j]−price[i]}
//		sell[i][j] = max{sell[i−1][j],buy[i−1][j−1]+price[i]}
func maxProfit_2(k int, prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}

	k = min(k, n/2)
	buy := make([][]int, n)
	sell := make([][]int, n)
	for i := 0; i < n; i++ {
		buy[i] = make([]int, k+1)
		sell[i] = make([]int, k+1)
	}
	// 设置初始值
	buy[0][0] = -prices[0]
	sell[0][0] = 0
	for i := 1; i < n; i++ {
		buy[i][0] = max(buy[i-1][0], sell[i-1][0]-prices[i])
		sell[i][0] = 0
	}
	for j := 1; j < k+1; j++ {
		buy[0][j] = math.MinInt32 / 2
		sell[0][j] = math.MinInt32 / 2
	}

	// 状态转移
	for i := 1; i < n; i++ {
		for j := 1; j < k+1; j++ {
			buy[i][j] = max(buy[i-1][j], sell[i-1][j]-prices[i])
			sell[i][j] = max(sell[i-1][j], buy[i-1][j-1]+prices[i])
		}
	}
	return max(sell[n-1]...)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(a ...int) int {
	n := len(a)
	if n == 0 {
		return 0
	}

	ans := a[0]
	for i := 1; i < n; i++ {
		if a[i] > ans {
			ans = a[i]
		}
	}
	return ans
}
