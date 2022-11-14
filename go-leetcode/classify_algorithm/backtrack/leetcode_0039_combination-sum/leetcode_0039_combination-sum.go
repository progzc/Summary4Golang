package leetcode_0039_combination_sum

// 0039. 组合总和
// https://leetcode.cn/problems/combination-sum/

// combinationSum dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n)
// 特点:
//	a.无重复元素
//	b.同一个数字可以无限制重复被选择
func combinationSum(candidates []int, target int) [][]int {
	var (
		ans  [][]int
		comb []int
		n    = len(candidates)
		dfs  func(begin, target int)
	)

	// begin: 表示搜索起点
	// target: 表示目标值
	dfs = func(begin, target int) {
		if target < 0 {
			return
		}

		if target == 0 {
			ans = append(ans, append([]int(nil), comb...))
			return
		}

		// 重点理解这里从 begin 开始搜索的语意
		for i := begin; i < n; i++ {
			comb = append(comb, candidates[i])
			// 注意：由于每一个元素可以重复使用，下一轮搜索的起点依然是 i，这里非常容易弄错
			dfs(i, target-candidates[i])
			comb = comb[:len(comb)-1]
		}
	}
	dfs(0, target)
	return ans
}

// combinationSum_2 dfs
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(n)
// 特点:
//	a.无重复元素
//	b.同一个数字可以无限制重复被选择
func combinationSum_2(candidates []int, target int) [][]int {
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
