package leetcode_0473_matchsticks_to_square

import "sort"

// 473. 火柴拼正方形
// https://leetcode.cn/problems/matchsticks-to-square/

// 与下面题目类似：
// 698. 划分为k个相等的子集
// https://leetcode.cn/problems/partition-to-k-equal-sum-subsets/

// makesquare dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(2^n)
func makesquare(matchsticks []int) bool {
	sum, edge := 0, 4
	for _, stick := range matchsticks {
		sum += stick
	}

	n := len(matchsticks)
	if n < edge || sum%edge != 0 {
		return false
	}

	// 降序排列：先让值大的元素选择桶，这样可以增加剪枝的命中率，从而降低回溯的概率
	sort.Slice(matchsticks, func(i, j int) bool {
		return matchsticks[i] > matchsticks[j]
	})

	var (
		dfs    func(idx int) bool
		bucket = make([]int, edge)
		target = sum / edge
	)
	dfs = func(idx int) bool {
		if idx == n {
			return true
		}
		for i := 0; i < edge; i++ {
			// 如果当前桶和上一个桶内的元素和相等，则跳过。
			// 原因：如果元素和相等，那么 matchsticks[idx] 选择上一个桶和选择当前桶可以得到的结果是一致的
			if i > 0 && bucket[i] == bucket[i-1] {
				continue
			}
			// 剪枝：放入球后超过 target 的值，选择一下桶
			if bucket[i]+matchsticks[idx] > target {
				continue
			}
			// 做选择：放入 i 号桶
			bucket[i] += matchsticks[idx]
			// 	处理下一个数
			if dfs(idx + 1) {
				return true
			}
			// 撤销选择：挪出 i 号桶
			bucket[i] -= matchsticks[idx]
		}
		return false
	}
	return dfs(0)
}
