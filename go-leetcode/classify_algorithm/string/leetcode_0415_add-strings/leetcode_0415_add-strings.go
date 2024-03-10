package leetcode_0415_add_strings

import (
	"fmt"
	"strconv"
)

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

// addStrings_2 模拟
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func addStrings_2(num1 string, num2 string) string {
	n1, n2 := len(num1)-1, len(num2)-1
	carry := 0
	var ans string
	for n1 >= 0 || n2 >= 0 {
		var x, y int
		if n1 >= 0 {
			x, _ = strconv.Atoi(string(num1[n1]))
			n1--
		}
		if n2 >= 0 {
			y, _ = strconv.Atoi(string(num2[n2]))
			n2--
		}
		tmp := x + y + carry
		cur := tmp % 10
		carry = tmp / 10
		ans = fmt.Sprintf("%d", cur) + ans
	}
	if carry > 0 {
		ans = fmt.Sprintf("%d", carry) + ans
	}
	return ans
}

// addStrings_3 模拟
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func addStrings_3(num1 string, num2 string) string {
	m, n := len(num1)-1, len(num2)-1
	remain := 0
	var ans string
	for m >= 0 || n >= 0 {
		x, y := 0, 0
		if m >= 0 {
			x = int(num1[m] - '0')
			m--
		}
		if n >= 0 {
			y = int(num2[n] - '0')
			n--
		}
		cur := (x + y + remain) % 10
		remain = (x + y + remain) / 10
		ans = strconv.Itoa(cur) + ans
	}
	if remain > 0 {
		ans = strconv.Itoa(remain) + ans
	}
	return ans
}
