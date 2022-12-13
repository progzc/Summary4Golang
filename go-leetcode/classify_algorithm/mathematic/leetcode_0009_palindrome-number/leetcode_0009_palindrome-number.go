package leetcode_0009_palindrome_number

import "fmt"

// 9. 回文数
// https://leetcode.cn/problems/palindrome-number/

// isPalindrome
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isPalindrome(x int) bool {
	str := fmt.Sprintf("%d", x)
	n := len(str)
	for i := 0; i < n/2; i++ {
		if str[i] != str[n-i-1] {
			return false
		}
	}
	return true
}

// isPalindrome_2
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func isPalindrome_2(x int) bool {
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}

	div := 1
	for x/div >= 10 {
		div *= 10
	}
	for x > 0 {
		l := x / div
		r := x % 10
		if l != r {
			return false
		}
		x = (x % div) / 10
		div /= 100
	}
	return true
}
