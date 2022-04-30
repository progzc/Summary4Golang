package leetcode_0067_add_binary

import (
	"strconv"
	"testing"
)

// 0067.二进制求和
// https://leetcode-cn.com/problems/add-binary/

func Test_addBinary(t *testing.T) {
	a, b := "11", "1"
	println(addBinary(a, b)) // 100
}

// addBinary
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func addBinary(a string, b string) string {
	m, n := len(a), len(b)
	maxLen := max(m, n)
	ans, carry := "", 0
	for i := 0; i < maxLen; i++ {
		if i < m {
			carry += int(a[m-i-1] - '0')
		}
		if i < n {
			carry += int(b[n-i-1] - '0')
		}
		// 注意下面两行的顺序
		ans = strconv.Itoa(carry%2) + ans
		carry /= 2
	}
	if carry > 0 {
		ans = "1" + ans
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
