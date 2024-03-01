package leetcode_1035_uncrossed_lines

// 1035. 不相交的线
// https://leetcode.cn/problems/uncrossed-lines/

// 类似题：
// 	72.编辑距离
// 	583.两个字符串的删除操作

// maxUncrossedLines
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：
//	动态规划公式：用dp[i][j]表示nums1的前i个数和nums2的前j个数之间可以绘制的最大连线树。则：
//		若nums1与nums2的最后一个数字相同: dp[i][j] = max{dp[i-1][j],dp[i][j-1],dp[i-1][j-1]+1}
//		若nums1与nums2的最后一个数字不同: dp[i][j] = max{dp[i-1][j],dp[i][j-1],dp[i-1][j-1]}
func maxUncrossedLines(nums1 []int, nums2 []int) int {
	n, m := len(nums1), len(nums2)
	if n == 0 || m == 0 {
		return 0
	}
	// dp数组及边界状态初始化
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	// 递推
	for i := 1; i < n+1; i++ {
		for j := 1; j < m+1; j++ {
			// 注意这里与【583.两个字符串的删除操作】的区别
			if nums1[i-1] == nums2[j-1] {
				dp[i][j] = max(dp[i-1][j], max(dp[i][j-1], dp[i-1][j-1]+1))
			} else {
				dp[i][j] = max(dp[i-1][j], max(dp[i][j-1], dp[i-1][j-1]))
			}
		}
	}
	return dp[n][m]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
