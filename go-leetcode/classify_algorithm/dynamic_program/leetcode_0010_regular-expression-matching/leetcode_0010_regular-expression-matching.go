package leetcode_0010_regular_expression_matching

// 10. 正则表达式匹配
// https://leetcode.cn/problems/regular-expression-matching/

// isMatch 动态规划
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路:
// 	状态: dp[i][j]表示s的前i个字符与p的前j个字符是否能够匹配。
//	转移方程:
//	  a.若p[j]是小写字母, 那么必须在s中匹配一个相同的小写字母,即:
//		当s[i] matches p[j], dp[i][j] = dp[i-1][j-1]
//		当s[i] not matches p[j], dp[i][j] = false
//	  b.若p[j]是'*', 那么就表示我们可以对 p 的第 j-1 个字符匹配任意自然数次。那么有下面两种情况:
//		i)匹配 s 末尾的一个字符，将该字符扔掉，而该组合还可继续进行匹配。
//		ii)不匹配字符，将该组合扔掉，不再进行匹配。
//		  当s[i] matches p[j-1], dp[i][j] = dp[i-1][j] || dp[i][j-2]，注意: 式子中, dp[i-1][j]表示可以继续匹配，dp[i][j-2]表示匹配0次。
//		  当s[i] not matches p[j-1], dp[i][j] = dp[i][j-2]
//	  c.若p[j]是'.', dp[i][j] = dp[i-1][j-1]
//	边界条件:
//		dp[0][0] = true，即两个空字符串是可以匹配的。
func isMatch(s string, p string) bool {
	var match func(x, y byte) bool
	match = func(sCh, pCh byte) bool {
		if pCh == '.' {
			return true
		}
		return sCh == pCh
	}

	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := 0; i < m+1; i++ {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for j := 1; j < n+1; j++ {
		if p[j-1] == '*' {
			dp[0][j] = dp[0][j-2]
		}
	}

	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			if p[j-1] == '*' {
				if match(s[i-1], p[j-2]) {
					dp[i][j] = dp[i-1][j] || dp[i][j-2]
				} else {
					dp[i][j] = dp[i][j-2]
				}
			} else if p[j-1] == '.' {
				dp[i][j] = dp[i-1][j-1]
			} else {
				if match(s[i-1], p[j-1]) {
					dp[i][j] = dp[i-1][j-1]
				} else {
					dp[i][j] = false
				}
			}
			//fmt.Printf("dp[%d][%d] = %v \n", i, j, dp[i][j])
		}
	}
	return dp[m][n]
}
