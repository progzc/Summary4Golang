package leetcode_0377_combination_sum_iv

// 0377. 组合总和 Ⅳ
// https://leetcode.cn/problems/combination-sum-iv/

// combinationSum4 超时
// 特点:
//	a.数组元素都是正数（假设若是负数呢?又该如何处理）
//	b.数组元素是不同的
//	c.数组元素可无限次取
//	d.顺序不同的序列被视作不同的组合
func combinationSum4(nums []int, target int) int {
	var (
		ans  int
		comb []int
		n    = len(nums)
		dfs  func(begin, target int)
	)

	// begin: 表示搜索起点
	// target: 表示目标值
	dfs = func(begin, target int) {
		if target < 0 {
			return
		}

		if target == 0 {
			ans++
			return
		}

		// 重点理解这里从 begin 开始搜索的语意
		for i := begin; i < n; i++ {
			comb = append(comb, nums[i])
			// 注意：由于每一个元素可以重复使用，下一轮搜索的起点依然是 begin，这里非常容易弄错
			dfs(begin, target-nums[i])
			comb = comb[:len(comb)-1]
		}
	}
	dfs(0, target)
	return ans
}

// combinationSum4_2 动态规划(组合背包)
// 时间复杂度: O(N*T)
// 空间复杂度: O(T)
// 思路：考虑顺序的组合背包
func combinationSum4_2(nums []int, target int) int {
	dp := make([]int, target+1)
	// 注意事项：只有当不选取任何元素时，元素之和才为0，因此只有1种方案
	dp[0] = 1
	for i := 1; i <= target; i++ {
		for j := 0; j < len(nums); j++ {
			if i >= nums[j] {
				dp[i] += dp[i-nums[j]]
			}
		}
	}
	return dp[target]
}
