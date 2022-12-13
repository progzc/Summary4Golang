package leetcode_0131_palindrome_partitioning

// 0131. 分割回文串
// https://leetcode.cn/problems/palindrome-partitioning/

// partition dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n^2)
// 注意: 结果集中不同的顺序也是不同的方案
//	针对用例:
//		实际输入: "aaa"
//		预期输出: [["a","a","a"],["a","aa"],["aa","a"],["aaa"]]
//		错误输出: [["a","a","a"],["a","aa"],["aaa"]]
func partition(s string) [][]string {
	var (
		ans          [][]string
		isPalindrome func(s string) bool
	)

	isPalindrome = func(s string) bool {
		sLen := len(s)
		for i := 0; i < sLen/2; i++ {
			if s[i] != s[sLen-i-1] {
				return false
			}
		}
		return true
	}

	n := len(s)
	for i := 1; i <= n; i++ {
		if isPalindrome(s[:i]) {
			if len(s[i:]) == 0 {
				ans = append(ans, []string{s[:i]})
			} else {
				nexts := partition(s[i:])
				for _, next := range nexts {
					ans = append(ans, append([]string{s[:i]}, next...))
				}
			}
		}
	}
	return ans
}

// partition_3 dfs+动态规划预处理
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n^2)
// 注意: 结果集中不同的顺序也是不同的方案
//	针对用例:
//		实际输入: "aaa"
//		预期输出: [["a","a","a"],["a","aa"],["aa","a"],["aaa"]]
//		错误输出: [["a","a","a"],["a","aa"],["aaa"]]
// 思路: 判断字符串是否为回文字符串可以使用动态规划进行预处理，可以加快速度。
//	状态: dp[i][j]表示s[i...j](左闭右闭区间)是否为回文串。
//	转移方程:
//		当j-i<0时, dp[i][j] = false
//		当j-i==0时,dp[i][j] = true
//		当j-i==1时,dp[i][j] = s[i] == s[j]
//		当j-i>=2时,dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
func partition_3(s string) [][]string {
	var (
		ans  [][]string
		comb []string
		dfs  func(idx int)
	)

	n := len(s)
	dp := make([][]bool, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]bool, n)
	}
	for i := n - 1; i >= 0; i-- {
		for j := i; j < n; j++ {
			if j-i == 0 {
				dp[i][j] = true
			} else if j-i == 1 {
				dp[i][j] = s[i] == s[j]
			} else {
				dp[i][j] = s[i] == s[j] && dp[i+1][j-1]
			}
		}
	}

	dfs = func(i int) {
		if i == n {
			ans = append(ans, append([]string(nil), comb...))
			return
		}
		for j := i; j < n; j++ {
			if dp[i][j] {
				comb = append(comb, s[i:j+1])
				dfs(j + 1)
				comb = comb[:len(comb)-1]
			}
		}
	}
	dfs(0)
	return ans
}
