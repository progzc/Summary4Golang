package leetcode_0055_jump_game

import "math"

// 0055. 跳跃游戏🌟
// https://leetcode.cn/problems/jump-game

// canJump 动态规划
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func canJump(nums []int) bool {
	n := len(nums)
	if n <= 1 {
		return true
	}
	dp := make([]bool, n)
	dp[0] = true
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if dp[j] && j+nums[j] >= i {
				dp[i] = true
				break
			}
		}
	}
	return dp[n-1]
}

// canJump_2 贪心算法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func canJump_2(nums []int) bool {
	n := len(nums)
	if n <= 1 {
		return true
	}
	i := 0
	for i < n-1 {
		if i+nums[i] >= n-1 {
			return true
		} else if nums[i] <= 0 {
			return false
		} else {
			cur := math.MinInt64
			k := i + 1
			for j := i + 1; j <= nums[i]+i; j++ {
				if j+nums[j] > cur {
					cur = j + nums[j]
					k = j
				}
			}
			i = k
		}
	}
	return false
}
