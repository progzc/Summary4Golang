package leetcode_0186_reverse_words_in_a_string_ii

// 0186.翻转字符串里的单词 II
// https://leetcode.cn/problems/reverse-words-in-a-string-ii/

// reverseWords
// 思路：两次翻转，第一次全局翻转，第二次每个单词进行翻转
// 时间复杂度：O(n)
// 空间复杂度：O(1)
func reverseWords(s []byte) {
	length := len(s)
	// 1.全局翻转
	reverse(s, 0, length-1)

	start := 0
	for i := 0; i < length; i++ {
		if s[i] == ' ' {
			// 2.1翻转前面的单词
			reverse(s, start, i-1)
			start = i + 1
		}
	}
	// 2.2 翻转最后一个单词
	reverse(s, start, length-1)
}

func reverse(s []byte, start, end int) {
	for start < end {
		s[start], s[end] = s[end], s[start]
		start++
		end--
	}
}
