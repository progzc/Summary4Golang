package futu_20221124_best_time_to_buy_and_sell_stock

// 题目:
// 给定一个int类型的数组, 表示一只股票最近N天的价格。
// 假设你每次买卖只能一股，可以买卖多次，但是手里最多只能持有一股。请写一个函数,
// 计算你所能获取的最大利润。
// 例如，一只股票最近N天的价格为[]int{1,4,2,3}，那么你所能获取的最大利润为4

// 力扣：
// 122. 买卖股票的最佳时机 II
// https://leetcode.cn/problems/best-time-to-buy-and-sell-stock-ii/

// maxProfit dfs 超时
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func maxProfit(nums []int) int {
	var (
		ans int
		dfs func(idx, pre, sum int, flag bool)
	)

	n := len(nums)
	dfs = func(idx, pre, sum int, flag bool) {
		if idx == n {
			return
		}
		// 未持有股票
		if !flag {
			// 选择买入
			dfs(idx+1, nums[idx], sum, true)
			// 选择不买入
			dfs(idx+1, pre, sum, false)
		} else {
			// 持有股票
			// 选择卖出
			ans = max(ans, sum+nums[idx]-pre)
			dfs(idx+1, 0, sum+nums[idx]-pre, false)
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
