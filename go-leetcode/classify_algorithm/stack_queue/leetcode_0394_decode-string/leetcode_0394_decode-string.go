package leetcode_0394_decode_string

import (
	"strconv"
	"strings"
)

// 0394. 字符串解码🌟
// https://leetcode.cn/problems/decode-string

// decodeString 栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func decodeString(s string) string {
	var stack []string
	for i := 0; i < len(s); {
		if s[i] >= '0' && s[i] <= '9' {
			j := i + 1
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			numStr := s[i:j]
			stack = append(stack, numStr)
			i = j
		} else if (s[i] >= 'a' && s[i] <= 'z') || s[i] == '[' {
			stack = append(stack, string(s[i]))
			i++
		} else {
			var sub string
			for stack[len(stack)-1] != "[" {
				sub = stack[len(stack)-1] + sub
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // 去掉"["
			num, _ := strconv.Atoi(stack[len(stack)-1])
			stack = stack[:len(stack)-1] // 去掉数字
			temp := strings.Repeat(sub, num)
			stack = append(stack, temp)
			i++
		}
	}
	return strings.Join(stack, "")
}
