package leetcode_0253_meeting_rooms_ii

import "sort"

// 0253.会议室 II
// https://leetcode-cn.com/problems/meeting-rooms-ii/

// minMeetingRooms 巧妙的思路
// 时间复杂度: O(nlog(n)) + O(n)
// 空间复杂度: O(nlog(n)) + O(1)
// 思路：所需会议室数量 等价于 要统计同一时刻进行的最大会议的数量。
//	a.把所有的开始时间和结束时间放在一起排序。
//	b.用cur表示当前进行的会议数量，遍历排序后的时间数组。
//	c.如果是开始时间，cur加1，如果是结束时间，cur减1。
//	d.在遍历的过程中，cur出现的最大值就是需要的房间数。
func minMeetingRooms(intervals [][]int) int {
	type item struct {
		time int // 会议时间
		flag int // 1表示开始会议,同时表示打开一个新的会议室；-1表示结束会议，同时表示关闭一个会议室
	}
	var ss []item
	for _, interval := range intervals {
		ss = append(ss, item{interval[0], 1})
		ss = append(ss, item{interval[1], -1})
	}

	// 注意： 输入：[[13,15],[1,13]]，预期输出：1；
	// 易错点：排序思想：先按照会议时间升序排列；当时间相同时，则先结束会议再开始会议
	sort.SliceStable(ss, func(i, j int) bool {
		return ss[i].time < ss[j].time || (ss[i].time == ss[j].time && ss[i].flag < ss[j].flag)
	})

	ans, cur := 0, 0
	for _, s := range ss {
		cur += s.flag
		ans = max(ans, cur)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
