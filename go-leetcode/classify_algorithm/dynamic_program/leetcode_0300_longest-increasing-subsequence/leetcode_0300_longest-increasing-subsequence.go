package leetcode_0300_longest_increasing_subsequence

// 300.最长递增子序列
// https://leetcode-cn.com/problems/longest-increasing-subsequence/

// lengthOfLIS 动态规划
// 时间复杂度：O(n*n)
// 空间复杂度：O(n)
// 状态：dp[i]：以nums[i]结尾的最长子序列的长度
// 转移方程：dp[i]=max(dp[j]+1),其中0<=j<i 且 nums[j]<nums[i]
func lengthOfLIS(nums []int) int {
	dp := make([]int, len(nums))
	dp[0] = 1
	ans := 1
	for i := 1; i < len(nums); i++ {
		dp[i] = 1
		for j := i - 1; j >= 0; j-- {
			if nums[j] < nums[i] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
		ans = max(ans, dp[i])
	}
	return ans
}

// lengthOfLIS_2 贪心+二分
// 时间复杂度：O(n*log(n))
// 空间复杂度：O(n)
// 状态：dp[i]：表示长度为i的最长上升子序列的末尾元素的最小值
// 可使用反证法证明dp[i]是关于i调递递增的；而根据dp[i]单调递增，可以使用二分法进行搜索
func lengthOfLIS_2(nums []int) int {
	// TODO
	return -1
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
