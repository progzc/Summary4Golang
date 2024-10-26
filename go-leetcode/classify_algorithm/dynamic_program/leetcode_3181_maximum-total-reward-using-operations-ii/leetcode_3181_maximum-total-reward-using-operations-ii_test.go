package leetcode_3181_maximum_total_reward_using_operations_ii

import (
	"sort"
)

// 3181. 执行操作可获得的最大总奖励 II
// https://leetcode.cn/problems/maximum-total-reward-using-operations-ii

// 注意：与【3180. 执行操作可获得的最大总奖励 I】题目一模一样，区别在于数据规模不一样

// maxTotalReward 动态规划（0/1背包）会超时
// 时间复杂度: O(n(m+log(n)))
// 空间复杂度: O(m+log(n))
func maxTotalReward(rewardValues []int) int {
	// 假设上一次操作选择的奖励值为 x1, 那么执行操作后的总奖励 x ≥ x1，
	// 根据题意，后面任一操作选择的奖励值 x2 一定都大于 x，从而有 x2 > x1。
	// 因此执行的操作是按照奖励值单调递增的。
	sort.Ints(rewardValues)
	// 记 rewardValues 的最大值为 m，因为最后一次操作前的总奖励一定小于等于 m−1，所以可获得的最大总奖励小于等于 2m−1。
	m := rewardValues[len(rewardValues)-1]
	// dp[i]表示总奖励 i 是否可以获得
	dp := make([]bool, 2*m)
	// dp[0]=true 表示不执行任何操作获得总奖励 0
	dp[0] = true
	for _, x := range rewardValues {
		for k := 2*x - 1; k >= x; k-- {
			if dp[k-x] {
				dp[k] = true
			}
		}
	}

	ans := 0
	for i := len(dp) - 1; i >= 0; i-- {
		if dp[i] {
			ans = i
			break
		}
	}
	return ans
}
