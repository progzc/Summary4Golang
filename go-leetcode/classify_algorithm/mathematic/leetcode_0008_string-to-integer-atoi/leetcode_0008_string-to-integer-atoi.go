package leetcode_0008_string_to_integer_atoi

import "math"

// 8. 字符串转换整数 (atoi)
// https://leetcode.cn/problems/string-to-integer-atoi/

// myAtoi
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func myAtoi(s string) int {
	abs, sign, i, n := 0, 1, 0, len(s)
	// 去掉前导空格
	for i < n && s[i] == ' ' {
		i++
	}
	// 标记正负号
	if i < n {
		if s[i] == '-' {
			sign = -1
			i++
		} else if s[i] == '+' {
			sign = 1
			i++
		}
	}
	for ; i < n && s[i] >= '0' && s[i] <= '9'; i++ {
		abs = 10*abs + int(s[i]-'0')
		if sign*abs < math.MinInt32 {
			return math.MinInt32
		} else if sign*abs > math.MaxInt32 {
			return math.MaxInt32
		}
	}
	return sign * abs
}
