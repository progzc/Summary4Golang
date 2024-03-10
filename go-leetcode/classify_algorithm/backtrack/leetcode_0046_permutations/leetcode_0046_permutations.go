package leetcode_0046_permutations

// 0046.全排列
// https://leetcode-cn.com/problems/permutations/
// 题意：给定一个 不含重复 数字的数组nums，返回其所有可能的全排列。你可以按任意顺序返回答案。

// permute 深度优先遍历
// 时间复杂度: O(n*n!)
// 空间复杂度: O(n)
// 优点：可以借助原数组来标记已选择哪些数
// 缺点：这样生成的全排列不是按照字典序存储的
func permute(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}
	var (
		ans    [][]int
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
		for i := idx; i < len(nums); i++ {
			// 选择：选择第i个位置作为第idx个数
			output[idx], output[i] = output[i], output[idx]
			// 递归：填下一个位置
			dfs(idx+1, output)
			// 回溯：在下一次选择之前，必须回撤销上一次的选择
			output[idx], output[i] = output[i], output[idx]
		}
	}

	output = append([]int(nil), nums...)
	dfs(0, output)
	return ans
}

// permute_2 深度优先遍历
// 时间复杂度: O(n*n!)
// 空间复杂度: O(2n)
// 优点：这样生成的全排列是按照字典序存储的
// 缺点：需要耗费额外的空间来标记已选择哪些数
func permute_2(nums []int) [][]int {
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
			// 剪枝：只可以选择未选择过的数
			if used[i] {
				continue
			}
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
	dfs(0, output)
	return ans
}
