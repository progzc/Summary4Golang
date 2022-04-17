package dynamic_program

import (
	"fmt"
	"testing"
)

// 此类题型总结：https://leetcode-cn.com/problems/last-stone-weight-ii/solution/yi-pian-wen-zhang-chi-tou-bei-bao-wen-ti-5lfv/

// Test_knapsack_01 01背包问题
func Test_knapsack_01(t *testing.T) {
	weight, values, volume := []int{1, 2, 3, 4}, []int{2, 4, 4, 5}, 5
	fmt.Println(knapsack_01(weight, values, volume)) // 8
}

// Test_knapsack_01_optimize 01背包问题(优化)
func Test_knapsack_01_optimize(t *testing.T) {
	weight, values, volume := []int{1, 2, 3, 4}, []int{2, 4, 4, 5}, 5
	fmt.Println(knapsack_01_optimize(weight, values, volume)) // 8
}

// knapsack_01 01背包问题
// 时间复杂度：O(W*V)
// 空间复杂度：O(W*V)
// 思路：
//	二维状态：dp[i][j]表示当背包容量为j时,在第0~i个物品中选择物品的最大价值
// 输入：
//	weight 各个物品的重量
//	values 各个物品对应的价值
//	volume 背包的容量
// 输出：背包装下物品的最大价值
func knapsack_01(weights, values []int, volume int) int {
	// 定义二维dp
	dp := make([][]int, len(weights))
	for i := range dp {
		dp[i] = make([]int, volume+1)
	}
	// 放入第一个物品
	for j := 0; j <= volume; j++ {
		if weights[0] <= j {
			dp[0][j] = weights[0]
		}
	}

	// 初始化二维dp
	for i := 1; i < len(weights); i++ {
		for j := 0; j <= volume; j++ {
			if weights[i] <= j {
				// 背包容量足够拿第i个物品,可拿可不拿
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-weights[i]]+values[i])
			} else {
				// 容量不足以拿第i个物品
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[len(weights)-1][volume]
}

// knapsack_01_optimize 01背包问题(优化)
// 时间复杂度：O(W*V)
// 空间复杂度：O(V)
// 思路：
//	一维状态：dp[j]表示当背包容量为j时,选择物品的最大价值
// 输入：
//	weight 各个物品的重量
//	values 各个物品对应的价值
//	volume 背包的容量
// 输出：背包装下物品的最大价值
func knapsack_01_optimize(weights, values []int, volume int) int {
	// 定义二维dp
	dp := make([]int, volume+1)

	// 初始化二维dp
	for i := 0; i < len(weights); i++ {
		// 注意这里一定要逆序，仔细思考下原因
		for j := volume; j >= 0; j-- {
			if weights[i] <= j {
				// 背包容量足够拿第i个物品,可拿可不拿
				dp[j] = max(dp[j], dp[j-weights[i]]+values[i])
			}
		}
	}
	return dp[volume]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
