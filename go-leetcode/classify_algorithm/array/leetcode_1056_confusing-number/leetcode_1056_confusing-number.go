package leetcode_1056_confusing_number

// 1056. 易混淆数
// https://leetcode.cn/problems/confusing-number

// confusingNumber 顺序法
// 时间复杂度：O(n)
// 空间复杂度: O(1)
func confusingNumber(n int) bool {
	m := map[int]int{0: 0, 1: 1, 6: 9, 8: 8, 9: 6}
	t := n
	x := 0
	for n > 0 {
		mod := n % 10
		n = n / 10
		if v, ok := m[mod]; !ok {
			return false
		} else {
			x = x*10 + v
		}
	}
	return t != x
}
