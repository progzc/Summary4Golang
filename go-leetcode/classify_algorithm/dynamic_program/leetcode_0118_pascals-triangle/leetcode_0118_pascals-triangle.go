package leetcode_0118_pascals_triangle

// 0118. 杨辉三角
// https://leetcode.cn/problems/pascals-triangle

// generate 动态规划
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func generate(numRows int) [][]int {
	var (
		ans [][]int
		pre []int
	)
	if numRows < 0 {
		return ans
	}
	if numRows >= 1 {
		pre = []int{1}
	}
	ans = append(ans, pre)
	for i := 1; i < numRows; i++ {
		cur := make([]int, i+1)
		for j := 0; j < len(cur)-1; j++ {
			if j-1 < 0 {
				cur[j] = pre[j]
			} else {
				cur[j] = pre[j-1] + pre[j]
			}
		}
		cur[len(cur)-1] = 1
		ans = append(ans, cur)
		pre = cur
	}
	return ans
}
