package leetcode_0122_best_time_to_buy_and_sell_stock_ii

import "math"

// 0122.买卖股票的最佳时机 II
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

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
