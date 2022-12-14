package leetcode_0006_zigzag_conversion

import "strings"

// 6. Z 字形变换
// https://leetcode.cn/problems/zigzag-conversion/

// convert 模拟
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func convert(s string, numRows int) string {
	if numRows < 2 {
		return s
	}
	n := len(s)
	rows := make([]strings.Builder, numRows)
	pos, flag := 0, -1
	for i := 0; i < n; i++ {
		rows[pos].WriteByte(s[i])
		if pos == 0 || pos == numRows-1 {
			flag = -flag
		}
		pos += flag
	}
	var ans string
	for _, row := range rows {
		ans += row.String()
	}
	return ans
}
