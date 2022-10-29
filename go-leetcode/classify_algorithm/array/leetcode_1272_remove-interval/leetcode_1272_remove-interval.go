package leetcode_1272_remove_interval

// 1272. 删除区间
// https://leetcode.cn/problems/remove-interval/

// removeInterval 分类讨论
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func removeInterval(intervals [][]int, toBeRemoved []int) [][]int {
	var ans [][]int
	toL, toR := toBeRemoved[0], toBeRemoved[1]
	for _, interval := range intervals {
		x, y := interval[0], interval[1]
		if toL >= y || toR <= x {
			ans = append(ans, []int{x, y})
		} else {
			if toL > x {
				ans = append(ans, []int{x, toL})
			}
			if toR < y {
				ans = append(ans, []int{toR, y})
			}
		}
	}
	return ans
}
