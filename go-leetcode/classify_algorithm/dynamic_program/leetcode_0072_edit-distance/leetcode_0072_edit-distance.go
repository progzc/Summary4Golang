package leetcode_0072_edit_distance

import "math"

// 0072.编辑距离🌟
// https://leetcode-cn.com/problems/edit-distance/

// 注意:
// 【583.两个字符串的删除操作】：只允许插入
// 【72.编辑距离】：允许删除/插入/修改

// minDistance
// 时间复杂度：O(m*n)
// 空间复杂度：O(m*n)
// 思路：
//
//	本质不同的操作只有三种：
//		a.在单词A中插入一个字符(等价于在单词B中删除一个字符)
//		b.在单词B中插入一个字符(等价于在单词A中删除一个字符)
//		c.修改单词A的一个字符
//	动态规划公式：用D[i][j]表示A的前i个字母和B的前j个字母之间的编辑距离。则
//		若A和B的最后一个字母不同：D[i][j] = min{ D[i-1][j]+1, D[i][j-1]+1, D[i-1][j-1]+1 }
//		若A和B的最后一个字母相同：D[i][j] = min{ D[i-1][j]+1, D[i][j-1]+1, D[i-1][j-1] }
//
// 扩展：编辑距离算法被数据科学家广泛应用，是用作机器翻译和语音识别评价标准的基本算法。
func minDistance(word1 string, word2 string) int {
	n, m := len(word1), len(word2)
	// 若有一个字符串为空串
	if n == 0 || m == 0 {
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
	// 递推
	for i := 1; i < n+1; i++ {
		for j := 1; j < m+1; j++ {
			// 注意这里与【583.两个字符串的删除操作】的区别
			if word1[i-1] == word2[j-1] {
				dp[i][j] = min(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1])
			} else {
				dp[i][j] = min(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+1)
			}
		}
	}
	return dp[n][m]
}

// minDistance_2
// 时间复杂度：O(m*n)
// 空间复杂度：O(m*n)
// 思路：
//
//	本质不同的操作只有三种：
//		a.在单词A中插入一个字符(等价于在单词B中删除一个字符)
//		b.在单词B中插入一个字符(等价于在单词A中删除一个字符)
//		c.修改单词A的一个字符
//	动态规划公式：用D[i][j]表示A的前i个字母和B的前j个字母之间的编辑距离。则
//		若A和B的最后一个字母不同：D[i][j] = min{ D[i-1][j]+1, D[i][j-1]+1, D[i-1][j-1]+1 }
//		若A和B的最后一个字母相同：D[i][j] = min{ D[i-1][j]+1, D[i][j-1]+1, D[i-1][j-1] }
//
// 扩展：编辑距离算法被数据科学家广泛应用，是用作机器翻译和语音识别评价标准的基本算法。
func minDistance_2(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		if i == 0 {
			if word1[i] != word2[0] {
				dp[i][0] = 1
			}
		} else {
			if word1[i] == word2[0] {
				dp[i][0] = min(dp[i-1][0]+1, i)
			} else {
				dp[i][0] = min(dp[i-1][0]+1, i+1)
			}
		}
	}

	for j := 0; j < n; j++ {
		if j == 0 {
			if word1[0] != word2[j] {
				dp[0][j] = 1
			}
		} else {
			if word1[0] == word2[j] {
				dp[0][j] = min(dp[0][j-1]+1, j)
			} else {
				dp[0][j] = min(dp[0][j-1]+1, j+1)
			}
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if word1[i] == word2[j] {
				dp[i][j] = min(dp[i-1][j-1], min(dp[i-1][j]+1, dp[i][j-1]+1))
			} else {
				dp[i][j] = min(dp[i-1][j-1]+1, min(dp[i-1][j]+1, dp[i][j-1]+1))
			}
		}
	}
	return dp[m-1][n-1]
}

func min(nums ...int) int {
	if len(nums) == 0 {
		return math.MinInt64
	}
	minNum := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < minNum {
			minNum = nums[i]
		}
	}
	return minNum
}
