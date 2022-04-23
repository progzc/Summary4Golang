package leetcode_0005_longest_palindromic_substring

// 0005.最长回文子串
// https://leetcode-cn.com/problems/longest-palindromic-substring/

// longestPalindrome 动态规划
// 时间复杂度: O(n*n)
// 空间复杂度: O(n*n)
// 状态：dp[i][j]表示字符串s的下标范围[i,j]内的字符串是否是回文
// 边界条件：
//	a.当0<=i<=j<n，才会有dp[i][j]=true，否则dp[i][j]=false
//	b.对于任意的0<=i<n,都有dp[i][i]=true
//	c.
// 转移方程：
//	当i+1<=j-1时，dp[i][j]=dp[i+1]dp[j-1] && s[i]==s[j]
//	当i+1>j-1时，dp[i][j]=s[i]==s[j]
func longestPalindrome(s string) string {
	n := len(s)
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	maxLen, start, end := 1, 0, 0
	// i必须逆序
	// 思考：如果i顺序会怎么样?有些结果还未计算出来
	for i := n - 1; i >= 0; i-- {
		dp[i][i] = true
		for j := i + 1; j < n; j++ {
			// 如果i+1>j-1，即j-i+1<3, 即len(s[i,j])<3, 则
			// 这个条件不能掉，否则针对"cbdd"用例会出错
			if i+1 > j-1 {
				dp[i][j] = s[i] == s[j]
			} else {
				dp[i][j] = dp[i+1][j-1] && s[i] == s[j]
			}
			if dp[i][j] && j-i+1 > maxLen {
				maxLen = j - i + 1
				start, end = i, j
			}
		}
	}
	return s[start : end+1]
}
