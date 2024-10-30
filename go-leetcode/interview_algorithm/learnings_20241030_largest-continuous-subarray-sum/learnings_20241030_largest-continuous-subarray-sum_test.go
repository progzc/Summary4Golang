package learnings_20241030_largest_continuous_subarray_sum

import (
	"fmt"
	"testing"
)

// 乐信圣文一面
// 寻找最大的连续子数组和
func TestLargestSubarraySum(t *testing.T) {
	nums := []int{1, 2, 3, -4, -5, 10, 11, -4, 10}
	fmt.Println(largestSubarraySum(nums)) // 27

	nums = []int{1, 2, 3, -4, -5, 10, 11}
	fmt.Println(largestSubarraySum(nums)) // 21
}

// largestSubarraySum 动态规划
func largestSubarraySum(nums []int) int {
	var ans int
	n := len(nums)
	if n == 0 {
		return ans
	}
	dp := make([]int, n)
	dp[0] = nums[0]

	for i := 1; i < n; i++ {
		if dp[i-1] > 0 {
			dp[i] = dp[i-1] + nums[i]
		} else {
			dp[i] = nums[i]
		}
		if dp[i] > ans {
			ans = dp[i]
		}
	}
	return ans
}
