package leetcode_0301_remove_invalid_parentheses

// 301. 删除无效的括号
// https://leetcode.cn/problems/remove-invalid-parentheses/

// removeInvalidParentheses dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n^2)
// 思路:
//	a.通过一次遍历计算出各自最少应该删除的左括号和右括号的数量, 具体步骤如下:
//	  i)当遍历到「左括号」的时候，「左括号」数量+1。
//	  ii)当遍历到「右括号」的时候，
//		 如果此时「左括号」的数量不为 0，因为「右括号」可以与之前遍历到的「左括号」匹配，此时「左括号」出现的次数 -1。
//		 如果此时「左括号」的数量为 00，「右括号」数量+1。
//	b.尝试在原字符串 s 中去掉 lremove 个左括号和 rremove 个右括号，然后检测剩余的字符串是否合法匹配，如果合法匹配则我们则认为该字符串为可能的结果，
// 	  我们利用回溯算法来尝试搜索所有可能的去除括号的方案。
//	  剪枝:
//		i)从字符串中每去掉一个括号，则更新 lremove 或者 rremove，当我们发现剩余未尝试的字符串的长度小于 lremove + rremove 时，则停止本次搜索。
//		ii)当 lremove 和 rremove 同时为 0 时，则我们检测当前的字符串是否合法匹配，如果合法匹配则我们将其记录下来。
//	  去重:
//		i)可以使用哈希表去重
//		ii)我们在每次进行搜索时，如果遇到连续相同的括号我们只需要搜索一次即可，比如当前遇到的字符串为 "(((())"，去掉前四个左括号中的任意一个，
//		   生成的字符串是一样的，均为 "((())"，因此我们在尝试搜索时，只需去掉一个左括号进行下一轮搜索，不需要将前四个左括号都尝试一遍。
func removeInvalidParentheses(s string) []string {
	lmove, rmove := 0, 0
	for _, ch := range s {
		if ch == '(' {
			lmove++
		} else if ch == ')' {
			if lmove == 0 {
				rmove++
			} else {
				lmove--
			}
		}
	}

	var (
		isValid func(str string) bool
		dfs     func(str string, start, lmove, rmove int)
		check   = make(map[string]bool)
		ans     []string
	)
	// isValid 判断字符串str是否是有效的括号
	isValid = func(str string) bool {
		cnt := 0
		for _, ch := range str {
			if ch == '(' {
				cnt++
			} else if ch == ')' {
				cnt--
				if cnt < 0 {
					return false
				}
			}
		}
		return cnt == 0
	}

	dfs = func(str string, start, lmove, rmove int) {
		if lmove < 0 || rmove < 0 {
			return
		}
		if lmove == 0 && rmove == 0 {
			if isValid(str) && !check[str] {
				ans = append(ans, str)
				check[str] = true
			}
			return
		}
		for i := start; i < len(str); i++ {
			if i > start && str[i] == str[i-1] {
				continue
			}
			// 如果剩余字符无法满足去掉的数量要求，直接返回
			// 下面这个条件等价于 lmove+rmove > (len(str)-1)-i+1
			if lmove+rmove > len(str)-i {
				return
			}
			// 尝试去掉第一个左括号
			if lmove > 0 && str[i] == '(' {
				dfs(str[:i]+str[i+1:], i, lmove-1, rmove)
			}
			// 尝试去掉一个右括号
			if rmove > 0 && str[i] == ')' {
				dfs(str[:i]+str[i+1:], i, lmove, rmove-1)
			}
		}
	}
	dfs(s, 0, lmove, rmove)
	return ans
}
