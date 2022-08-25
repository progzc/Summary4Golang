package leetcode_0254_factor_combinations

import "math"

// 0254. 因子的组合
// https://leetcode.cn/problems/factor-combinations/
// 注意事项：
// a.主要考虑怎么避免重复
// b.输出结果的顺序是一定的

// getFactors 递归
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
// 思路：
//	1.dfs递归: 遍历范围为i=[2,sqrt(n)+1)的因子
//	2.单层递归: 如果i是n的因子,该层中添加一个[]int{i,n/i}的组合
//	3.递归处理:
//		3.1 查看n/i能否从因子i开始被拆分成多因子组合的子序列
//		3.2 如果子序列存在, 则在该层结果中继续添加一个[]int{i,子序列}的新组合
func getFactors(n int) [][]int {
	var dfs func(n, from int) [][]int

	dfs = func(n, from int) [][]int {
		var ans [][]int
		to := int(math.Sqrt(float64(n))) + 1
		// 注意这里不能是：i <= to
		for i := from; i < to; i++ {
			if n%i == 0 {
				ans = append(ans, []int{i, n / i})
				// 注意这里递归是：dfs(n/i, i)
				for _, k := range dfs(n/i, i) {
					ans = append(ans, append([]int{i}, k...))
				}
			}
		}
		return ans
	}
	return dfs(n, 2)
}
