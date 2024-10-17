package leetcode_0322_coin_change

import "math"

// 0322.零钱兑换
// https://leetcode-cn.com/problems/coin-change/

// 关于从0-1背包到完全背包的推导过程如下:
// https://leetcode.cn/problems/coin-change/solution/by-flix-su7s/

// coinChange_3 动态规划（完全背包的最值问题）二维动态规划
// 时间复杂度：O(S*n)
// 空间复杂度：O(S)
// 思路：
//
//		状态定义: dp[i][j] 表示：从前 i 种硬币中组成金额 j 所需最少的硬币数量。
//	 转移方程:
//			优化前: dp[i][j] = min{ dp[i−1][j], dp[i−1][j−k⋅wi]+k} ,0<=k⋅wi<=j,其中k=1,2,3...
//			优化后: dp[i][j] = min{ dp[i−1][j], dp[i][j−wi]+1} ,0<=k⋅wi<=j
func coinChange_3(coins []int, amount int) int {
	n := len(coins)
	dp := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]int, amount+1)
	}

	// 初始化
	// 初始化时，不合法的或未定义的状态则可以设置为正无穷或一个不可能取到的较大值
	// dp[0][0]=0：表示从前 0 种硬币中选出若干个组成金额 0 所对应的最小硬币数目为 0，即「空集合」不选任何硬币即可得到金额 0。
	// 对于其他 dp[0][j],j≥1，则可将其设置为正无穷或一个不可能取到的较大值，例如 dp[0][j] = +INF：「空集合」中无法选出任何硬币组成金额 j≥1。
	for j := 0; j < amount+1; j++ {
		if j == 0 {
			dp[0][0] = 0
		} else {
			dp[0][j] = math.MaxInt32
		}
	}
	for i := 1; i < n+1; i++ {
		for j := 0; j < amount+1; j++ {
			if coins[i-1] > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = min(dp[i-1][j], dp[i][j-coins[i-1]]+1)
			}
		}
	}
	ans := dp[n][amount]
	if ans != math.MaxInt32 {
		return ans
	}
	return -1
}

// coinChange_4 动态规划（完全背包的最值问题）一维动态规划
// 时间复杂度：O(S*n)
// 空间复杂度：O(S)
// 思路：
//
//		状态定义: dp[i][j] 表示：从前 i 种硬币中组成金额 j 所需最少的硬币数量。
//	 转移方程:
//			优化前: dp[i][j] = min{ dp[i−1][j], dp[i−1][j−k⋅wi]+k} ,0<=k⋅wi<=j,其中k=1,2,3...
//			优化后: dp[i][j] = min{ dp[i−1][j], dp[i][j−wi]+1} ,0<=k⋅wi<=j
func coinChange_4(coins []int, amount int) int {
	dp := make([]int, amount+1)
	// 初始化
	// 初始化时，不合法的或未定义的状态则可以设置为正无穷或一个不可能取到的较大值
	// dp[0][0]=0：表示从前 0 种硬币中选出若干个组成金额 0 所对应的最小硬币数目为 0，即「空集合」不选任何硬币即可得到金额 0。
	// 对于其他 dp[0][j],j≥1，则可将其设置为正无穷或一个不可能取到的较大值，例如 dp[0][j] = +INF：「空集合」中无法选出任何硬币组成金额 j≥1。
	for j := 0; j < amount+1; j++ {
		if j == 0 {
			dp[j] = 0
		} else {
			dp[j] = math.MaxInt32
		}
	}
	for _, coin := range coins {
		for j := 0; j < amount+1; j++ {
			if j >= coin {
				dp[j] = min(dp[j], dp[j-coin]+1)
			}
		}
	}
	ans := dp[amount]
	if ans != math.MaxInt32 {
		return ans
	}
	return -1
}

// coinChange 递归
// 时间复杂度：O(S*n)
// 空间复杂度：O(S)
// 思路：假设我们知道F(S)，即组成金额S最少的硬币数，最后一枚硬币的面值是C。
//
//	那么由于问题的最优子结构，转移方程应为：F(S)=F(S-C)+1
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
