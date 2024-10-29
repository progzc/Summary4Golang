package leetcode_3211_generate_binary_strings_without_adjacent_zeros

import (
	"fmt"
	"testing"
)

// 3211. 生成不含相邻零的二进制字符串
// https://leetcode.cn/problems/generate-binary-strings-without-adjacent-zeros

func TestValidStrings(t *testing.T) {
	fmt.Println(validStrings(3)) // [010 011 101 110 111]
	fmt.Println(validStrings(1)) // [0 1]
}

// validStrings 回溯
func validStrings(n int) []string {
	var (
		ans []string
		dfs func(i int, s string)
	)

	dfs = func(i int, s string) {
		if i == n {
			ans = append(ans, s)
			return
		}
		if len(s) > 0 && s[len(s)-1] == '0' {
			dfs(i+1, s+"1") // 选择 0
		} else {
			dfs(i+1, s+"0") // 选择 0
			dfs(i+1, s+"1") // 选择 1
		}
	}
	dfs(0, "")
	return ans
}
