package leetcode_0125_valid_palindrome

import "strings"

// 0125.验证回文串
// https://leetcode-cn.com/problems/valid-palindrome/

// isPalindrome 筛选+判断
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func isPalindrome(s string) bool {
	var ss string
	for i := 0; i < len(s); i++ {
		if isNumOrLetter(s[i]) {
			ss += string(s[i])
		}
	}
	n := len(ss)
	ss = strings.ToLower(ss)
	for i := 0; i < n/2; i++ {
		if ss[i] != ss[n-1-i] {
			return false
		}
	}
	return true
}

func isNumOrLetter(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}

// isPalindrome_2 双指针+调用库转换大小写
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：先转换为小写 + 使用双指针 （使用字符串"abcdef"这一特例来检验）
func isPalindrome_2(s string) bool {
	// 注意: 可以直接调用库转换为小写
	s = strings.ToLower(s)
	left, right := 0, len(s)-1
	for left < right {
		for left < right && !isNumOrLetter(s[left]) {
			left++
		}
		for left < right && !isNumOrLetter(s[right]) {
			right--
		}
		if left < right {
			if s[left] != s[right] {
				return false
			}
			left++
			right--
		}
	}
	return true
}

// isPalindrome_3 双指针+自己转换大小写
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：不使用库转换为小写 + 使用双指针 （使用字符串"abcdef"这一特例来检验）
func isPalindrome_3(s string) bool {
	// 注意: 不直接调用库转换为小写
	var ss string
	for i := 0; i < len(s); i++ {
		if isUpperLetter(s[i]) {
			ss = ss + string(s[i]+'a'-'A')
		} else {
			ss = ss + string(s[i])
		}
	}

	left, right := 0, len(ss)-1
	for left < right {
		for left < right && !isNumOrLetter(ss[left]) {
			left++
		}
		for left < right && !isNumOrLetter(ss[right]) {
			right--
		}
		if left < right {
			if ss[left] != ss[right] {
				return false
			}
			left++
			right--
		}
	}
	return true
}

func isUpperLetter(ch byte) bool {
	return ch >= 'A' && ch <= 'Z'
}
