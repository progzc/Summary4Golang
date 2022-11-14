package leetcode_0040_combination_sum_ii

import "sort"

// 0040. 组合总和 II
// https://leetcode.cn/problems/combination-sum-ii/

// combinationSum2
// 特点:
//	a.含重复元素
//	b.每个元素只能使用一次
func combinationSum2(candidates []int, target int) [][]int {
	// 0 排序
	sort.Ints(candidates)

	var ans [][]int
	var comb []int
	// idx 从候选数组的idx位置开始搜索
	// target 表示剩余
	var dfs func(idx, target int)
	dfs = func(idx, target int) {
		// 2. 添加结果
		if target == 0 {
			ans = append(ans, append([]int(nil), comb...))
			return
		}
		for i := idx; i < len(candidates); i++ {
			// 1.1 终止条件（大剪枝）
			// 大剪枝：减去 candidates[i] 小于 0，减去后面的 candidates[i + 1]、candidates[i + 2] 肯定也小于 0，因此用 break
			if target-candidates[i] < 0 {
				break
			}
			// 1.2 终止条件（小剪枝）
			// 小剪枝：同一层相同数值的结点，从第 2 个开始，候选数更少，结果一定发生重复，因此跳过，用 continue
			if i > idx && candidates[i] == candidates[i-1] {
				continue
			}
			// 3.1 选择当前数
			comb = append(comb, candidates[i])
			// 3.2 递归
			dfs(i+1, target-candidates[i])
			// 3.3 回退
			comb = comb[:len(comb)-1]
		}
	}
	dfs(0, target)
	return ans
}
