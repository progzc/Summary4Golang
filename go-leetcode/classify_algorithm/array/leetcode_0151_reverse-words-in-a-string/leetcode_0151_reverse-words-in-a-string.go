package leetcode_0151_reverse_words_in_a_string

import "strings"

// 0151. 反转字符串中的单词
// https://leetcode.cn/problems/reverse-words-in-a-string/

// reverseWords 标准库
// 思路：两次翻转，第一次全局翻转，第二次每个单词进行翻转
//	时间复杂度: O(n)
//	空间复杂度: O(n)
func reverseWords(s string) string {
	ss := strings.Fields(s)
	l := len(ss)
	for i := 0; i < l/2; i++ {
		ss[i], ss[l-1-i] = ss[l-1-i], ss[i]
	}
	return strings.Join(ss, " ")
}

// reverseWords_2 指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func reverseWords_2(s string) string {
	// a.使用双指针删除冗余的空格
	slow, fast := 0, 0
	b := []byte(s)
	// 删除头部冗余空格
	for fast < len(b) && b[fast] == ' ' {
		fast++
	}
	// 删除中间冗余空格
	for ; fast < len(b); fast++ {
		if fast-1 > 0 && b[fast-1] == b[fast] && b[fast] == ' ' {
			continue
		}
		b[slow] = b[fast]
		slow++
	}
	// 删除尾部冗余空格
	if slow-1 > 0 && b[slow-1] == ' ' {
		b = b[:slow-1]
	} else {
		b = b[:slow]
	}
	// b.反转整个字符串
	reverse(b, 0, len(b)-1)
	// c.二次反转
	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == ' ' {
			reverse(b, start, i-1)
			start = i + 1
		}
	}
	// 翻转最后一个单词
	reverse(b, start, len(b)-1)
	return string(b)
}

func reverse(s []byte, start, end int) {
	for start < end {
		s[start], s[end] = s[end], s[start]
		start++
		end--
	}
}
