package leetcode_0712_minimum_ascii_delete_sum_for_two_strings

// 712. 两个字符串的最小ASCII删除和
// https://leetcode.cn/problems/minimum-ascii-delete-sum-for-two-strings/

// 与下面类似：
// 【583.两个字符串的删除操作】

// minimumDeleteSum
// 时间复杂度：O(m*n)
// 空间复杂度：O(m*n)
// 思路:
//	本质不同的操作只有三种：
//		a.在单词A中插入一个字符(等价于在单词B中删除一个字符)
//		b.在单词B中插入一个字符(等价于在单词A中删除一个字符)
//	动态规划公式：用D[i][j]表示A的前i个字母和B的前j个字母相同的最小ASCII删除和，则
//		若A和B的最后一个字母不同：D[i][j] = min{ D[i-1][j] + s1[i-1], D[i][j-1] + s2[j-1] }
//		若A和B的最后一个字母相同：D[i][j] = D[i-1][j-1]
func minimumDeleteSum(s1 string, s2 string) int {
	n, m := len(s1), len(s2)
	dp := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]int, m+1)
	}

	for i := 1; i < n+1; i++ {
		dp[i][0] = dp[i-1][0] + int(s1[i-1])
	}
	for j := 1; j < m+1; j++ {
		dp[0][j] = dp[0][j-1] + int(s2[j-1])
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if s1[i] == s2[j] {
				dp[i+1][j+1] = dp[i][j]
			} else {
				dp[i+1][j+1] = min(dp[i][j+1]+int(s1[i]), dp[i+1][j]+int(s2[j]))
			}
		}
	}
	return dp[n][m]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
