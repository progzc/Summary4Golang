package leetcode_1087_brace_expansion

import (
	"sort"
	"strings"
)

// 1087. 花括号展开
// https://leetcode.cn/problems/brace-expansion/

// expand dfs
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
// 注意：expand_2比expand效率更高
func expand(s string) []string {
	if strings.Index(s, "{") == -1 {
		return []string{s}
	}

	var (
		dfs func(cur int, path string)
		ans []string
	)
	dfs = func(cur int, path string) {
		if cur == len(s) {
			ans = append(ans, path)
			return
		}
		if s[cur] != '{' {
			dfs(cur+1, path+string(s[cur]))
		} else {
			i := cur
			for s[i] != '}' {
				i++
			}
			elems := strings.Split(s[cur+1:i], ",")
			sort.Slice(elems, func(i, j int) bool {
				return elems[i] < elems[j]
			})
			for _, elem := range elems {
				dfs(i+1, path+elem)
			}
		}
	}

	dfs(0, "")
	return ans
}

// expand_2 dfs
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
func expand_2(s string) []string {
	if strings.Index(s, "{") == -1 {
		return []string{s}
	}

	var (
		dfs func(s string, b []byte, idx int)
		ans []string
	)
	dfs = func(s string, b []byte, idx int) {
		if len(s) == idx {
			ans = append(ans, string(b))
			return
		}

		if s[idx] == '{' {
			// 先计算出{}中内容的长度
			count := 0
			for i := idx + 1; s[i] != '}'; i++ {
				count++
			}
			// 下次要跳转的位置为idx+count+2
			for i := idx + 1; s[i] != '}'; i++ {
				if s[i] != ',' {
					b = append(b, s[i])
					dfs(s, b, idx+count+2)
					b = b[:len(b)-1]
				}
			}
		} else {
			b = append(b, s[idx])
			dfs(s, b, idx+1)
		}
	}
	dfs(s, nil, 0)
	sort.Slice(ans, func(i, j int) bool {
		return ans[i] < ans[j]
	})
	return ans
}
