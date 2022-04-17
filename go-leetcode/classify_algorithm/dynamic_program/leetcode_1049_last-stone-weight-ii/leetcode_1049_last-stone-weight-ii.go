package leetcode_1049_last_stone_weight_ii

// 1049.最后一块石头的重量 II
// https://leetcode-cn.com/problems/last-stone-weight-ii/

// lastStoneWeightII 动态规划
// 时间复杂度：
// 空间复杂度：
// 思路：把一堆石头分成两堆，求两堆石头重量差最小值；进一步地，要让差值小,两堆石头的重量都要接近sum/2。
//	设：石头总重量为sum，k=-1的石头的重量之和为neg，则k=1的石头的重量之和为sum-neg
//	那么：最后一块石头的重量=(sum-neg)-neg=sum-2*neg，要使最后一块石头的重量尽可能地小，neg 需要在不超过 sum/2 的前提下尽可能地大。
//	因此本问题可以看作是背包容量为 sum/2，物品重量和价值均为 stones(i)的0-1背包问题。
func lastStoneWeightII(stones []int) int {
	// 计算石头总重量
	sum := 0
	for _, stone := range stones {
		sum += stone
	}
	// 求解01背包问题
	target := sum / 2
	dp := make([]int, target+1)
	for _, stone := range stones {
		for i := target; i >= 0; i-- {
			if i >= stone {
				dp[i] = max(dp[i], dp[i-stone]+stone)
			}
		}
	}
	return sum - 2*dp[target]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
