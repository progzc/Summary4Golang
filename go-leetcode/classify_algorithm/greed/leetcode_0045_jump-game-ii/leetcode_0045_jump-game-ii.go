package leetcode_0045_jump_game_ii

import "math"

// 0045. 跳跃游戏 II
// https://leetcode.cn/problems/jump-game-ii

// jump 贪心算法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func jump(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return 0
	}
	count := 0
	i := 0
	for i < n-1 {
		if i+nums[i] >= n-1 {
			count++
			return count
		} else if nums[i] <= 0 {
			break
		} else {
			cur := math.MinInt64
			k := i + 1
			for j := i + 1; j <= nums[i]+i; j++ {
				if j+nums[j] > cur {
					cur = j + nums[j]
					k = j
				}
			}
			count++
			i = k
		}
	}
	return count
}

// jump_2 动态规划
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
func jump_2(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return 0
	}

	type Item struct {
		Arrived bool
		Cnt     int
	}
	dp := make([]*Item, n)
	dp[0] = &Item{true, 0}
	for i := 1; i < n; i++ {
		dp[i] = &Item{Arrived: false, Cnt: math.MaxInt64}
		for j := 0; j < i; j++ {
			if dp[j].Arrived && j+nums[j] >= i {
				dp[i].Arrived = true
				dp[i].Cnt = min(dp[i].Cnt, dp[j].Cnt+1)
			}
		}
	}
	return dp[n-1].Cnt
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
