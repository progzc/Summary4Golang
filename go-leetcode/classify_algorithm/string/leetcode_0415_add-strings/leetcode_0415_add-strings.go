package leetcode_0415_add_strings

import "strconv"

// 0415.字符串相加
// https://leetcode-cn.com/problems/add-strings/

// addStrings 模拟
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func addStrings(num1 string, num2 string) string {
	carry := 0
	ans := ""

	for i, j := len(num1)-1, len(num2)-1; i >= 0 || j >= 0; i, j = i-1, j-1 {
		var x, y int
		if i >= 0 {
			x = int(num1[i] - '0')
		}
		if j >= 0 {
			y = int(num2[j] - '0')
		}
		result := x + y + carry
		carry = result / 10
		ans = strconv.Itoa(result%10) + ans
	}
	if carry > 0 {
		ans = strconv.Itoa(carry) + ans
	}
	return ans
}
