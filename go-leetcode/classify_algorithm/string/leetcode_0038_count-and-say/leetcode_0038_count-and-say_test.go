package leetcode_0038_count_and_say

import (
	"strconv"
	"testing"
)

// 0038.外观数列
// https://leetcode.cn/problems/count-and-say/

func Test_countAndSay(t *testing.T) {
	println(countAndSay(10))
}

// countAndSay 遍历
// 时间复杂度: O(n*m)
// 空间复杂度: O(m)
func countAndSay(n int) string {
	if n == 1 {
		return "1"
	}

	pre := "1"
	for i := 1; i < n; i++ {
		ans := ""
		for j := 0; j < len(pre); {
			k := j + 1
			for k < len(pre) && pre[k-1] == pre[k] {
				k++
			}
			ans += strconv.Itoa(len(pre[j:k])) + string(pre[j])
			j = k
		}
		pre = ans
	}

	return pre
}
