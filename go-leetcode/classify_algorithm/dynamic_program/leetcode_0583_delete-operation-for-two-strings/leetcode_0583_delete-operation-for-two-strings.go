package leetcode_0583_delete_operation_for_two_strings

// 583. 两个字符串的删除操作
// https://leetcode.cn/problems/delete-operation-for-two-strings/

// 注意:
// 【583.两个字符串的删除操作】：只允许插入
// 【72.编辑距离】：允许删除/插入/修改

// 也可以借鉴：
// 【1143.最长公共子序列】

// minDistance
// 时间复杂度：O(m*n)
// 空间复杂度：O(m*n)
// 思路：
//	本质不同的操作只有三种：
//		a.在单词A中插入一个字符(等价于在单词B中删除一个字符)
//		b.在单词B中插入一个字符(等价于在单词A中删除一个字符)
//	动态规划公式：用D[i][j]表示A的前i个字母和B的前j个字母之间的最小步数。则
//		若A和B的最后一个字母不同：D[i][j] = min{ D[i-1][j]+1, D[i][j-1]+1, D[i-1][j-1]+2 }
//		若A和B的最后一个字母相同：D[i][j] = min{ D[i-1][j]+1, D[i][j-1]+1, D[i-1][j-1] }
func minDistance(word1 string, word2 string) int {
	n, m := len(word1), len(word2)
	// 若有一个字符串为空串
	if n*m == 0 {
		return n + m
	}

	// dp数组及边界状态初始化
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 0; i < n+1; i++ {
		dp[i][0] = i
	}
	for j := 0; j < m+1; j++ {
		dp[0][j] = j
	}

	for i := 1; i < n+1; i++ {
		for j := 1; j < m+1; j++ {
			s1 := dp[i-1][j] + 1
			s2 := dp[i][j-1] + 1
			// 注意这里与【72.编辑距离】的区别
			s3 := dp[i-1][j-1] + 2
			if word1[i-1] == word2[j-1] {
				s3 -= 2
			}
			dp[i][j] = min(s1, min(s2, s3))
		}
	}
	return dp[n][m]
}

// minDistance_2
// 时间复杂度：O(m*n)
// 空间复杂度：O(m*n)
// 思路:
//	可以计算两个字符串的最长公共子序列的长度，然后分别计算两个字符串的长度和最长公共子序列的长度之差，即为最小操作步数。
func minDistance_2(word1 string, word2 string) int {
	n, m := len(word1), len(word2)
	// 若有一个字符串为空串
	if n*m == 0 {
		return n + m
	}

	// dp数组及边界状态初始化
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if word1[i] == word2[j] {
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}
	return n + m - dp[n][m]*2
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
