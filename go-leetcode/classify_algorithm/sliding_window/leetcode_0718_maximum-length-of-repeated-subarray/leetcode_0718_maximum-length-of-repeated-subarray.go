package leetcode_0718_maximum_length_of_repeated_subarray

// 0718. 最长重复子数组
// https://leetcode.cn/problems/maximum-length-of-repeated-subarray/

// findLength 动态规划
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
func findLength(nums1 []int, nums2 []int) int {
	m, n := len(nums1), len(nums2)
	if m < 1 || n < 1 {
		return 0
	}

	ans := 0
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
	}
	for i := 0; i < m; i++ {
		if nums1[i] == nums2[0] {
			dp[i][0] = 1
			ans = 1
		}
	}
	for j := 0; j < n; j++ {
		if nums1[0] == nums2[j] {
			dp[0][j] = 1
			ans = 1
		}
	}

	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if nums1[i] == nums2[j] {
				dp[i][j] = dp[i-1][j-1] + 1
			}
			ans = max(ans, dp[i][j])
		}
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// findLength_2 滑动窗口
// 时间复杂度: O()
// 空间复杂度: O()
func findLength_2(nums1 []int, nums2 []int) int {
	// TODO
	return 0
}
