package leetcode_3185_count_pairs_that_form_a_complete_day_ii

import (
	"fmt"
	"testing"
)

// 3185. 构成整天的下标对数目 II🌟
// https://leetcode.cn/problems/count-pairs-that-form-a-complete-day-ii

func TestCountCompleteDayPairs(t *testing.T) {
	hours := []int{12, 12, 30, 24, 24}
	fmt.Println(countCompleteDayPairs(hours))
}

// countCompleteDayPairs 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 题目有点类似于两数之和
// 思路: hours[i]+hours[j] 能够被 24 整除，只需 hours[i] 除以 24 的余数与 hours[j] 除以 24 的余数之和能够被 24 整除。
func countCompleteDayPairs(hours []int) int64 {
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
	return int64(ans)
}

// countCompleteDayPairs_2 动态规划（会超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
func countCompleteDayPairs_2(hours []int) int64 {
	n := len(hours)
	if n < 2 {
		return 0
	}

	pre := int64(0)
	for i := 1; i < n; i++ {
		cnt := int64(0)
		for j := 0; j < i; j++ {
			if (hours[j]+hours[i])%24 == 0 {
				cnt++
			}
		}
		pre = pre + cnt
	}
	return pre
}
