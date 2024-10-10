package qcloud_20241010_xiaoxiaole

import (
	"fmt"
	"strings"
	"testing"
)

func TestSolve(t *testing.T) {
	s := "abbbdddac"
	fmt.Println(solve(s)) // c

	s = "ac"
	fmt.Println(solve(s)) // ac
}

// 青云：消消乐，可以重复消除。
// 例如：输入 abbbdddac，输出 c
// 例如：输入 ac，输出 ac
// 注意：这道题区别于下面这道题。下面这道题是 2 个为一组进行消除，而本道题可以多个不限数量一起消除，难度大了很多
// 1047. 删除字符串中的所有相邻重复项
// https://leetcode.cn/problems/remove-all-adjacent-duplicates-in-string
func solve(s string) string {
	n := len(s)
	if n <= 1 {
		return s
	}

	var ans []string
	for i := 0; i < n; {
		if i < n-1 {
			if s[i] != s[i+1] {
				ans = append(ans, string(s[i]))
				i++
			} else {
				j := i
				for j < n && s[j] == s[i] {
					j++
				}
				if j == n {
					i = j + 1
				} else {
					i = j
				}
			}
		} else {
			ans = append(ans, string(s[i]))
			i++
		}
	}
	if strings.Join(ans, "") == s {
		return s
	}
	return solve(strings.Join(ans, ""))
}
