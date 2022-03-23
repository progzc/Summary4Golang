package leetcode_0121_best_time_to_buy_and_sell_stock

import "math"

// 0121.买卖股票的最佳时机
// https://leetcode-cn.com/problems/best-time-to-buy-and-sell-stock/

// maxProfit 贪心法
// 时间复杂度：O(n)
// 空间复杂度：O(1)
func maxProfit(prices []int) int {
	minPrice, maxProfile := math.MaxInt64, 0
	for _, price := range prices {
		if price < minPrice {
			minPrice = price
		}
		if price-minPrice > maxProfile {
			maxProfile = price - minPrice
		}
	}
	return maxProfile
}
