package leetcode_0435_non_overlapping_intervals

import (
	"sort"
)

// 435. 无重叠区间
// https://leetcode.cn/problems/non-overlapping-intervals/

// eraseOverlapIntervals 排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(log(n))
func eraseOverlapIntervals(intervals [][]int) int {
	if len(intervals) <= 1 {
		return 0
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0] ||
			(intervals[i][0] == intervals[j][0] && intervals[i][1] < intervals[j][1])
	})

	count := 0
	pre := intervals[0]
	for i := 1; i < len(intervals); i++ {
		cur := intervals[i]
		if cur[0] >= pre[1] {
			pre = cur
		} else {
			if cur[1] <= pre[1] {
				pre = cur
			}
			count++
		}
	}
	return count
}

// eraseOverlapIntervals_2 动态规划（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
// 思路:
//	状态: dp[i]表示以区间 i 为最后一个区间，可以选出的区间数量的最大值
func eraseOverlapIntervals_2(intervals [][]int) int {
	if len(intervals) <= 1 {
		return 0
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	n := len(intervals)
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		dp[i] = 1
	}

	maxSelect := 1
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if intervals[j][1] <= intervals[i][0] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
		maxSelect = max(maxSelect, dp[i])
	}
	return n - maxSelect
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
