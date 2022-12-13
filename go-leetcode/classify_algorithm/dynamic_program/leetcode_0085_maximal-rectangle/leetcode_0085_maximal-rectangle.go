package leetcode_0085_maximal_rectangle

// 85. 最大矩形
// https://leetcode.cn/problems/maximal-rectangle/

// 单调栈解法类似于下面:
// 84. 柱状图中最大的矩形
// https://leetcode.cn/problems/largest-rectangle-in-histogram/

// maximalRectangle 暴力法
// 时间复杂度: O(m^2*n)
// 空间复杂度: O(m*n)
// 思路:
//	第一步: 先统计每个点结尾的横向最长连续1的个数
//	第二步: 以该点向上搜索, 以最小值为长, 搜索高度为宽, 计算面积
func maximalRectangle(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	// width 保存以当前数字结尾的连续 1 的个数
	width := make([][]int, m)
	for i := 0; i < m; i++ {
		width[i] = make([]int, n)
	}

	maxArea := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			// 更新width
			if matrix[i][j] == '1' {
				if j == 0 {
					width[i][j] = 1
				} else {
					width[i][j] = width[i][j-1] + 1
				}
			} else {
				width[i][j] = 0
			}
			// 记录所有行中最小的数
			minWidth := width[i][j]
			// 向上扩展行
			for k := i; k >= 0 && width[k][j] != '0'; k-- {
				height := i - k + 1
				minWidth = min(minWidth, width[k][j])
				maxArea = max(maxArea, minWidth*height)
			}
		}
	}
	return maxArea
}

// maximalRectangle_2 单调栈
// 时间复杂度: O(m^2*n)
// 空间复杂度: O(m*n)
// 思路: 同【84. 柱状图中最大的矩形】
//	求出每一层的heights[]然后传给【84. 柱状图中最大的矩形】的函数求解即可,求解过程中记录最大值。
func maximalRectangle_2(matrix [][]byte) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])
	maxArea := 0
	heights := make([]int, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] == '1' {
				heights[j] += 1
			} else {
				heights[j] = 0
			}
		}
		maxArea = max(maxArea, largestRectangleArea(heights))
	}
	return maxArea
}

// largestRectangleArea 单调栈+哨兵（空间换时间）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：单调递增(非严格)栈
// 注意：此函数即为【84. 柱状图中最大的矩形】的解法
func largestRectangleArea(heights []int) int {
	n := len(heights)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return heights[0]
	}

	// 首尾加入哨兵
	heights = append([]int{0}, heights...)
	heights = append(heights, 0)

	var (
		stack []int
		ans   int
	)
	for i, h := range heights {
		for len(stack) > 0 && h < heights[stack[len(stack)-1]] {
			curH := heights[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			curW := i - stack[len(stack)-1] - 1
			ans = max(ans, curH*curW)
		}
		stack = append(stack, i)
	}
	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
