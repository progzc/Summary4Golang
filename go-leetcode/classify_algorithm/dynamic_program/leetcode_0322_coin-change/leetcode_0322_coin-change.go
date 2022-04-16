package leetcode_0322_coin_change

import "math"

// 0322.零钱兑换
// https://leetcode-cn.com/problems/coin-change/

// coinChange 递归
// 时间复杂度：O(S*n)
// 空间复杂度：O(S)
// 思路：假设我们知道F(S)，即组成金额S最少的硬币数，最后一枚硬币的面值是C。
//		那么由于问题的最优子结构，转移方程应为：F(S)=F(S-C)+1
func coinChange(coins []int, amount int) int {
	if amount < 1 {
		return 0
	}
	count := make([]int, amount)
	var change func(rem int) int
	change = func(rem int) int {
		if rem < 0 {
			return -1
		}
		if rem == 0 {
			return 0
		}
		if count[rem-1] != 0 {
			return count[rem-1]
		}
		minV := math.MaxInt32
		for _, coin := range coins {
			res := change(rem - coin)
			if res >= 0 && res < minV {
				minV = 1 + res
			}
		}
		if minV == math.MaxInt32 {
			count[rem-1] = -1
		} else {
			count[rem-1] = minV
		}
		return count[rem-1]
	}
	return change(amount)
}

// coinChange_2 动态规划
// 时间复杂度：O(S*n)
// 空间复杂度：O(S)
// 思路：将递归写成迭代的方式
func coinChange_2(coins []int, amount int) int {
	max := amount + 1
	dp := make([]int, amount+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = max
	}

	dp[0] = 0
	for i := 1; i <= amount; i++ {
		for j := 0; j < len(coins); j++ {
			if coins[j] <= i {
				dp[i] = min(dp[i], dp[i-coins[j]]+1)
			}
		}
	}
	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
