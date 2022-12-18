package leetcode_0012_integer_to_roman

import "strings"

// 12. 整数转罗马数字
// https://leetcode.cn/problems/integer-to-roman/

// intToRoman 模拟
// 时间复杂度: O(1)
// 空间复杂度: O(1)
// 思路: 核心思想是贪心
func intToRoman(num int) string {
	var (
		str func(num *int) string
		ans string
	)
	str = func(num *int) string {
		if *num >= 1000 {
			*num -= 1000
			return "M"
		} else if *num >= 900 {
			*num -= 900
			return "CM"
		} else if *num >= 500 {
			*num -= 500
			return "D"
		} else if *num >= 400 {
			*num -= 400
			return "CD"
		} else if *num >= 100 {
			*num -= 100
			return "C"
		} else if *num >= 90 {
			*num -= 90
			return "XC"
		} else if *num >= 50 {
			*num -= 50
			return "L"
		} else if *num >= 40 {
			*num -= 40
			return "XL"
		} else if *num >= 10 {
			*num -= 10
			return "X"
		} else if *num >= 9 {
			*num -= 9
			return "IX"
		} else if *num >= 5 {
			*num -= 5
			return "V"
		} else if *num >= 4 {
			*num -= 4
			return "IV"
		} else if *num >= 1 {
			*num -= 1
			return "I"
		}
		return ""
	}
	for num > 0 {
		ans += str(&num)
	}
	return ans
}

// intToRoman_2 模拟（代码优化）
// 时间复杂度: O(1)
// 空间复杂度: O(1)
// 思路: 核心思想是贪心
func intToRoman_2(num int) string {
	var (
		sb     strings.Builder
		values = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
		romans = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	)
	for i := 0; i < len(values); i++ {
		for num >= values[i] {
			num -= values[i]
			sb.WriteString(romans[i])
		}
		if num == 0 {
			break
		}
	}
	return sb.String()
}
