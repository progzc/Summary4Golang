package leetcode_0216_combination_sum_iii

// 0216. 组合总和 III
// https://leetcode.cn/problems/combination-sum-iii/

// combinationSum3
// 特点:
// 	a.只使用数字1到9
//	b.每个数字最多使用一次
func combinationSum3(k int, n int) [][]int {
	var (
		ans  [][]int
		comb []int
		dfs  func(begin, count, target int)
	)

	start, end := 1, 9
	// begin 从数字几开始进行选择
	// count 已经选择了多少个数
	// target 剩下数的和
	dfs = func(begin, count, target int) {
		if count == k && target == 0 {
			ans = append(ans, append([]int(nil), comb...))
			return
		}

		for i := begin; i <= end; i++ {
			if target-i < 0 {
				break
			}
			comb = append(comb, i)
			dfs(i+1, count+1, target-i)
			comb = comb[:len(comb)-1]
		}
	}
	dfs(start, 0, n)
	return ans
}
