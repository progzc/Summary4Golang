package leetcode_0123_best_time_to_buy_and_sell_stock_iii

// 123. 买卖股票的最佳时机 III
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-iii/

// maxProfit dfs（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func maxProfit(prices []int) int {
	var (
		ans int
		dfs func(idx, count, pre, sum int, flag bool)
	)

	n := len(prices)
	k := 2
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
// 空间复杂度: O(1)
// 思路:
//	状态: 由于我们最多可以完成两笔交易，因此在任意一天结束之后，我们会处于以下五个状态中的一种：
//		未进行过任何操作；
//		buy1：只进行过一次买操作；
//		sell1：进行了一次买操作和一次卖操作，即完成了一笔交易；
//		buy2：在完成了一笔交易的前提下，进行了第二次买操作；
//		sell2：完成了全部两笔交易。
//		由于第一个状态的利润显然为 0，因此我们可以不用将其记录。对于剩下的四个状态，我们分别将它们的最大利润记为 buy1,sell1,buy2,sell2
//	转移方程:
//		buy1 = max(buy1’, -prices[i])
//		sell1 = max(sell1‘, buy1’+prices[i])
//		buy2 = max(buy2‘, sell1’-prices[i])
//		sell2 = max(sell2‘, buy2’+prices[i])
//	注意事项:
//		相当于有4种独立的状态。
func maxProfit_2(prices []int) int {
	n := len(prices)
	if n < 2 {
		return 0
	}
	// 初始状态:
	// buy1即以prices[0]的价格买入股票，所以：buy1 = -prices[0]
	// sell1即为在同一天买入并卖出，所以sell1 = 0
	// buy2即为在同一天买入并且卖出后再以 prices[0] 的价格买入股票，所以buy2 = -prices[0]
	// sell2即为即为在同一天买入和卖出，并再次买入和卖出，所以sell2 = 0
	buy1, sell1 := -prices[0], 0
	buy2, sell2 := -prices[0], 0
	for i := 1; i < n; i++ {
		oldBuy1, oldSell1, oldBuy2, oldSell2 := buy1, sell1, buy2, sell2
		buy1 = max(oldBuy1, -prices[i])
		sell1 = max(oldSell1, oldBuy1+prices[i])
		buy2 = max(oldBuy2, oldSell1-prices[i])
		sell2 = max(oldSell2, oldBuy2+prices[i])
	}
	return max(max(0, sell1), sell2)
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
