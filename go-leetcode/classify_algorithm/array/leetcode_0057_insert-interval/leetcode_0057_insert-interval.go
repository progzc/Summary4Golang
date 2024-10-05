package leetcode_0057_insert_interval

// 0057. 插入区间
// https://leetcode.cn/problems/insert-interval

func insert(intervals [][]int, newInterval []int) [][]int {
	var ans [][]int
	if len(intervals) == 0 {
		ans = append(ans, newInterval)
		return ans
	}
	pre := newInterval
	for i := 0; i < len(intervals); i++ {
		if pre[1] < intervals[i][0] {
			ans = append(ans, pre)
			ans = append(ans, intervals[i:]...)
			break
		} else if pre[1] <= intervals[i][1] {
			ans = append(ans, []int{min(pre[0], intervals[i][0]), max(pre[1], intervals[i][1])})
			if i+1 < len(intervals) {
				ans = append(ans, intervals[i+1:]...)
			}
			break
		} else if pre[0] > intervals[i][1] {
			ans = append(ans, intervals[i])
			if i == len(intervals)-1 {
				ans = append(ans, pre)
			}
		} else {
			pre = []int{min(pre[0], intervals[i][0]), max(pre[1], intervals[i][1])}
			if i == len(intervals)-1 {
				ans = append(ans, pre)
			}
		}
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
