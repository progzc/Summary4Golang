package leetcode_0252_meeting_rooms

import "sort"

// 0252. 会议室
// https://leetcode.cn/problems/meeting-rooms/

// canAttendMeetings 排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(1)
func canAttendMeetings(intervals [][]int) bool {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	for i := 0; i < len(intervals)-1; i++ {
		if intervals[i][1] > intervals[i+1][0] {
			return false
		}
	}
	return true
}
