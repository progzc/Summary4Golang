package leetcode_0647_palindromic_substrings

// 0647.回文子串
// https://leetcode-cn.com/problems/palindromic-substrings/

// countSubstrings 动态规划
// 时间复杂度: O(n*n)
// 空间复杂度: O(n*n)
// 思路：与 0005.最长回文子串 思路相同
// 状态：dp[i][j]表示字符串s的下标范围[i,j]内的字符串是否是回文
// 边界条件：
//	a.当0<=i<=j<n，才会有dp[i][j]=true，否则dp[i][j]=false
//	b.对于任意的0<=i<n,都有dp[i][i]=true
//	c.
// 转移方程：
//	当i+1<=j-1时，dp[i][j]=dp[i+1]dp[j-1] && s[i]==s[j]
//	当i+1>j-1时，dp[i][j]=s[i]==s[j]
// 注意事项：具有不同开始位置或结束位置的子串，即使是由相同的字符组成，也会被视作不同的子串。
func countSubstrings(s string) int {
	n := len(s)
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	count := 0
	// i必须逆序
	// 思考：如果i顺序会怎么样?有些结果还未计算出来
	for i := n - 1; i >= 0; i-- {
		dp[i][i] = true
		if dp[i][i] {
			count++
		}

		for j := i + 1; j < n; j++ {
			// 如果i+1>j-1，即j-i+1<3, 即len(s[i,j])<3
			if i+1 > j-1 {
				dp[i][j] = s[i] == s[j]
			} else {
				dp[i][j] = dp[i+1][j-1] && s[i] == s[j]
			}
			if dp[i][j] {
				count++
			}
		}
	}
	return count
}
