package leetcode_0739_daily_temperatures

// 0739. 每日温度
// https://leetcode.cn/problems/daily-temperatures/

// dailyTemperatures 单调递减(非严格)栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func dailyTemperatures(temperatures []int) []int {
	var (
		stack []int
		ans   = make([]int, len(temperatures))
	)
	for i := 0; i < len(temperatures); i++ {
		for len(stack) > 0 && temperatures[stack[len(stack)-1]] < temperatures[i] {
			ans[stack[len(stack)-1]] = i - stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	return ans
}
