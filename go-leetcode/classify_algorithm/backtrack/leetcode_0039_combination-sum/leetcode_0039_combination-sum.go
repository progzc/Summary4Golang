package leetcode_0039_combination_sum

// 0039. 组合总和
// https://leetcode.cn/problems/combination-sum/

// combinationSum dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n)
func combinationSum(candidates []int, target int) [][]int {
	var (
		ans  [][]int
		comb []int
		dfs  func(i, target int)
	)

	n := len(candidates)
	dfs = func(i, target int) {
		if target < 0 {
			return
		}

		if target == 0 {
			ans = append(ans, append([]int(nil), comb...))
			return
		}

		// 没有元素可供做选择了
		if i >= n {
			return
		}

		// 选择当前元素
		comb = append(comb, candidates[i])
		dfs(i, target-candidates[i])
		comb = comb[:len(comb)-1]
		// 不选择当前元素
		dfs(i+1, target)
	}
	dfs(0, target)
	return ans
}
