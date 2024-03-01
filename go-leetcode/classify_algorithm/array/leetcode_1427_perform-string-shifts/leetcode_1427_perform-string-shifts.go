package leetcode_1427_perform_string_shifts

// 1427. 字符串的左右移
// https://leetcode.cn/problems/perform-string-shifts/

// stringShift 循环移动
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func stringShift(s string, shift [][]int) string {
	sum := 0
	direct := map[int]int{
		0: -1,
		1: 1,
	}
	for i := 0; i < len(shift); i++ {
		sum += direct[shift[i][0]] * shift[i][1]
	}
	n := len(s)
	mod := sum % n
	// 没有产生移动
	if mod == 0 {
		return s
	}
	// 左移统一转换成右移
	if mod < 0 {
		mod = mod + n
	}
	return s[n-mod:] + s[:n-mod]
}
