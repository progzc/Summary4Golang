package leetcode_0047_permutations_ii

// 0047.全排列 II
// https://leetcode-cn.com/problems/permutations-ii/
// 题意：给定一个 含重复 数字的数组nums，返回其所有可能的 不重复的 全排列。你可以按任意顺序返回答案。

func permuteUnique(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}
	var (
		ans    [][]int
		used   = make([]bool, len(nums))
		output []int
		dfs    func(idx int, output []int)
	)
	// dfs 表示从左往右填到第idx个位置,当前排列为output
	// 其中第0~idx-1的元素均已经填充过了，而idx~n-1的元素还未填充过
	dfs = func(idx int, output []int) {
		// 终止条件
		if idx == len(nums) {
			// 注意事项：ans = append(ans,output)这种写法是错误的
			ans = append(ans, append([]int(nil), output...))
		}
		for i := 0; i < len(nums); i++ {
			// 剪枝
			if used[i] || i > 0 && !used[i-1] && nums[i] == nums[i-1] {
				continue
			}

			// 前提：只可以选择未选择过的数
			if !used[i] {
				// 选择：选择第i个位置作为第idx个数
				output = append(output, nums[i])
				used[i] = true
				// 递归：填下一个位置
				dfs(idx+1, output)
				// 回溯：在下一次选择之前，必须回撤销上一次的选择
				used[i] = false
				output = output[:len(output)-1]
			}
		}
	}
	dfs(0, output)
	return ans
}
