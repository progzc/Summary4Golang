package leetcode_3184_count_pairs_that_form_a_complete_day_i

// 3184. 构成整天的下标对数目 I
// https://leetcode.cn/problems/count-pairs-that-form-a-complete-day-i

// countCompleteDayPairs 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 题目有点类似于两数之和
// 思路: hours[i]+hours[j] 能够被 24 整除，只需 hours[i] 除以 24 的余数与 hours[j] 除以 24 的余数之和能够被 24 整除。
func countCompleteDayPairs(hours []int) int {
	m := make(map[int]int)
	ans := 0
	for _, hour := range hours {
		if v, ok := m[(24-hour%24)%24]; ok {
			//下面这个容易错误，因为解决不了余数为 0 的情况
			//if v, ok := m[24-hour%24]; ok {
			ans += v
		}
		m[hour%24]++
	}
	return ans
}

// countCompleteDayPairs_2 动态规划
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func countCompleteDayPairs_2(hours []int) int {
	n := len(hours)
	if n < 2 {
		return 0
	}

	dp := make([]int, n)
	dp[0] = 0
	for i := 1; i < n; i++ {
		cnt := 0
		for j := 0; j < i; j++ {
			if (hours[j]+hours[i])%24 == 0 {
				cnt++
			}
		}
		dp[i] = dp[i-1] + cnt
	}
	return dp[n-1]
}
