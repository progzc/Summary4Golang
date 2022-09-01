package leetcode_1229_meeting_scheduler

import "sort"

// 01229. 安排会议日程
// https://leetcode.cn/problems/meeting-scheduler/

// minAvailableDuration 排序+双指针
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func minAvailableDuration(slots1 [][]int, slots2 [][]int, duration int) []int {
	// 先排序
	sort.Slice(slots1, func(i, j int) bool {
		return slots1[i][0] < slots1[j][0]
	})
	sort.Slice(slots2, func(i, j int) bool {
		return slots2[i][0] < slots2[j][0]
	})

	// 分情况讨论, 具体可参见 https://leetcode.cn/problems/meeting-scheduler/solution/golang-er-fen-cha-zhao-by-resara-7/
	for i, j := 0, 0; i < len(slots1) && j < len(slots2); {
		if slots1[i][0] >= slots2[j][1] {
			j++
		} else if slots1[i][1] >= slots2[j][1] && slots1[i][0] < slots2[j][1] {
			end := slots2[j][1]
			start := max(slots1[i][0], slots2[j][0])
			if end-start >= duration {
				return []int{start, start + duration}
			}
			j++
		} else if slots1[i][1] < slots2[j][1] && slots1[i][1] > slots2[j][0] {
			end := slots1[i][1]
			start := max(slots1[i][0], slots2[j][0])
			if end-start >= duration {
				return []int{start, start + duration}
			}
			i++
		} else {
			i++
		}
	}
	return []int{}
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
