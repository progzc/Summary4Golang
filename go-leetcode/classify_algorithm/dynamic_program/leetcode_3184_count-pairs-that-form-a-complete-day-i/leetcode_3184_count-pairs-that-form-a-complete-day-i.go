package leetcode_3184_count_pairs_that_form_a_complete_day_i

// 3184. 构成整天的下标对数目 I
// https://leetcode.cn/problems/count-pairs-that-form-a-complete-day-i

// countCompleteDayPairs 动态规划
func countCompleteDayPairs(hours []int) int {
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
