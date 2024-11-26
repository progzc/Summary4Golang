package leetcode_3206_alternating_groups_i

// 3206. 交替组 I
// https://leetcode.cn/problems/alternating-groups-i

// numberOfAlternatingGroups 顺序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func numberOfAlternatingGroups(colors []int) int {
	var ans int
	n := len(colors)
	if n <= 2 {
		return ans
	}
	for i := 0; i < n; i++ {
		need := (colors[i] + 1) % 2
		before := i - 1
		after := i + 1
		if before < 0 {
			before = n - 1
		}
		if after > n-1 {
			after = 0
		}
		if colors[before] == need && colors[after] == need {
			ans++
		}
	}
	return ans
}
