package leetcode_0007_reverse_integer

import "math"

// 0007.整数反转
// https://leetcode-cn.com/problems/reverse-integer/

// reverse
// 时间复杂度: O(20)
// 空间复杂度: O(1)
// 思路：枚举边界点的值
func reverse(x int) int {
	res := 0
	for x != 0 {
		digit := x % 10

		// res 是上一次反转后的值
		// 判断是否大于最大的32位整数
		if res > math.MaxInt32/10 || (res == math.MaxInt32/10 && digit > 7) {
			return 0
		}
		// 判断是否小于最小的32位整数
		if res < math.MinInt32/10 || (res == math.MinInt64/10 && digit < -8) {
			return 0
		}
		res = res*10 + digit
		x /= 10
	}
	return res
}
