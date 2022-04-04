package leetcode_0680_valid_palindrome_ii

// 0680.验证回文字符串Ⅱ
// https://leetcode-cn.com/problems/valid-palindrome-ii/

// validPalindrome 双指针+贪心
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// // 思路：遇到第一个不相同的字符，则若分别删除左右字符，剩下的字串至少有一个是回文串，则原字符串删除一个字符之后就以成为回文串。
// 对于用例："deeee",leetcode判断我的输出是false，实际在本地输出是true
var deleted = false

func validPalindrome(s string) bool {
	if len(s) == 0 {
		return false
	}
	left, right := 0, len(s)-1
	for left <= right {
		if s[left] == s[right] {
			left++
			right--
		} else {
			if !deleted {
				deleted = true
				return validPalindrome(s[left+1:right+1]) || validPalindrome(s[left:right])
			}
			return false
		}
	}
	return true
}

// validPalindrome_2 双指针+贪心
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：遇到第一个不相同的字符，则若分别删除左右字符，剩下的字串至少有一个是回文串，则原字符串删除一个字符之后就以成为回文串。
func validPalindrome_2(s string) bool {
	i, j := 0, len(s)-1
	for i <= j && s[i] == s[j] {
		i++
		j--
	}
	if i > j {
		return true
	} else {
		return palindrome(s[i+1:j+1]) || palindrome(s[i:j])
	}
}

func palindrome(s string) bool {
	i, j := 0, len(s)-1
	for i <= j && s[i] == s[j] {
		i++
		j--
	}
	return i > j
}
