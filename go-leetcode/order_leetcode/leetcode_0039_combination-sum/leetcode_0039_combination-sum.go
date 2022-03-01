package leetcode_0039_combination_sum

// 组合总和
// combinationSum 回溯法
// 时间复杂度 O(n*2^n)
// 空间复杂度 O(n)
func combinationSum(candidates []int, target int) [][]int {
	var comb []int
	var ans [][]int
	var dfs func(target, idx int)
	dfs = func(target, idx int) {
		// 1. 终止条件
		if idx == len(candidates) || target < 0 {
			return
		}
		// 2.记录结果
		if target == 0 {
			var temp []int
			temp = append(temp, comb...)
			ans = append(ans, temp)
			return
		}
		// 3. 做选择
		// 3.1 跳过当前位置的数
		dfs(target, idx+1)
		// 3.2 选择当前位置的数
		comb = append(comb, candidates[idx])
		dfs(target-candidates[idx], idx)
		// 4. 回退
		comb = comb[:len(comb)-1]
	}
	dfs(target, 0)
	return ans
}
