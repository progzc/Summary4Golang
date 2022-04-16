package leetcode_0072_edit_distance

// 0072.编辑距离
// https://leetcode-cn.com/problems/edit-distance/

// minDistance
// 时间复杂度：O(m*n)
// 空间复杂度：O(m*n)
// 思路：
//	本质不同的操作只有三种：
//		a.在单词A中插入一个字符(等价于在单词B中删除一个字符)
//		b.在单词B中插入一个字符(等价于在单词A中删除一个字符)
//		c.修改单词A的一个字符
//	动态规划公式：用D[i][j]表示A的前i个字母和B的前j个字母之间的编辑距离。则
//		若A和B的最后一个字母不同：D[i][j]=1+min{D[i-1][j],D[i][j-1],D[i-1][j-1]}
//		若A和B的最后一个字母相同：D[i][j]=1+min{D[i-1][j],D[i][j-1],D[i-1][j-1]-1}
// 扩展：编辑距离算法被数据科学家广泛应用，是用作机器翻译和语音识别评价标准的基本算法。
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
			s3 := dp[i-1][j-1] + 1
			if word1[i-1] == word2[j-1] {
				s3 -= 1
			}
			dp[i][j] = min(s1, min(s2, s3))
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
