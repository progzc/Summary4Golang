package leetcode_1208_get_equal_substrings_within_budget

// 1208. 尽可能使字符串相等
// https://leetcode.cn/problems/get-equal-substrings-within-budget/

// equalSubstring 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func equalSubstring(s string, t string, maxCost int) int {
	n := len(s)
	ans, cost := 0, 0
	for l, r := 0, 0; r < n; r++ {
		// 错误写法：cost += abs(int(s[r]-t[r]))，因为byte类型不能表示负数
		cost += abs(int(s[r]) - int(t[r]))
		for cost > maxCost {
			// 错误写法：cost -= abs(int(s[l]-t[l]))，因为byte类型不能表示负数
			cost -= abs(int(s[l]) - int(t[l]))
			l++
		}
		ans = max(ans, r-l+1)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
