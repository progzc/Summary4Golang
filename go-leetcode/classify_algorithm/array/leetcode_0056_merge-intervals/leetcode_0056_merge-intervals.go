package leetcode_0056_merge_intervals

import "sort"

// 0056.合并区间
// https://leetcode-cn.com/problems/merge-intervals/

// merge 排序
// 时间复杂度：O(nlog(n))
// 空间复杂度：O(log(n))
// 思路：可以使用反证法给出数学证明
//	a.将列表中的区间按照左端点升序排序。
//	b.然后我们将第一个区间加入 merged 数组中，并按顺序依次考虑之后的每个区间：
//		i)如果当前区间的左端点在数组 merged 中最后一个区间的右端点之后，那么它们不会重合，我们可以直接将这个区间加入数组 merged 的末尾;
//		ii)否则，它们重合，我们需要用当前区间的右端点更新数组 merged 中最后一个区间的右端点，将其置为二者的较大值。
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	// 排序是最关键的一步
	sort.Slice(intervals, func(i, j int) bool {
		// 按照区间的左端点排序
		return intervals[i][0] < intervals[j][0]
	})

	ans := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		r1 := ans[len(ans)-1][1]
		l2, r2 := intervals[i][0], intervals[i][1]
		// 下面这种情况一定不能合并
		if r1 < l2 {
			ans = append(ans, intervals[i])
		} else {
			ans[len(ans)-1][1] = max(r1, r2)
		}
	}
	return ans
}

// merge_2 排序
// 时间复杂度：O(nlog(n))
// 空间复杂度：O(log(n))
// 思路：可以使用反证法给出数学证明
//	a.将列表中的区间按照左端点升序排序。
//	b.然后我们将第一个区间加入 merged 数组中，并按顺序依次考虑之后的每个区间：
//		i)如果当前区间的左端点在数组 merged 中最后一个区间的右端点之后，那么它们不会重合，我们可以直接将这个区间加入数组 merged 的末尾;
//		ii)否则，它们重合，我们需要用当前区间的右端点更新数组 merged 中最后一个区间的右端点，将其置为二者的较大值。
func merge_2(intervals [][]int) [][]int {
	n := len(intervals)
	if n < 2 {
		return intervals
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0] || (intervals[i][0] == intervals[j][0] && intervals[i][1] < intervals[j][1])
	})

	var ans [][]int
	pre := intervals[0]
	for i := 1; i < n; i++ {
		cur := intervals[i]
		if cur[0] > pre[1] {
			ans = append(ans, pre)
			pre = cur
		} else {
			cur[0] = pre[0]
			cur[1] = max(cur[1], pre[1])
			pre = cur
		}
	}
	ans = append(ans, pre)
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
