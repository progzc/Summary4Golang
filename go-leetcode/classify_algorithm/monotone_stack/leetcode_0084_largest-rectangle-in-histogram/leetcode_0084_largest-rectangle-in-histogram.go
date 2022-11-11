package leetcode_0084_largest_rectangle_in_histogram

// 0084. 柱状图中最大的矩形
// https://leetcode.cn/problems/largest-rectangle-in-histogram/

// largestRectangleArea 暴力法（超时）
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
// 思路：枚举每个高度
func largestRectangleArea(heights []int) int {
	n := len(heights)
	ans := 0
	for i, h := range heights {
		l := i
		for l > 0 && heights[l-1] >= h {
			l--
		}

		r := i
		for r < n-1 && heights[r+1] >= h {
			r++
		}
		ans = max(ans, (r-l+1)*h)
	}
	return ans
}

// largestRectangleArea_2 单调栈+哨兵（空间换时间）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：单调递增(非严格)栈
func largestRectangleArea_2(heights []int) int {
	n := len(heights)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return heights[0]
	}

	// 首尾加入哨兵
	newHeights := make([]int, n+2)
	copy(newHeights[1:n+1], heights)

	var (
		stack []int
		ans   int
	)
	// 加入首部哨兵下标
	stack = append(stack, 0)
	for i := 1; i < n+2; i++ {
		for newHeights[i] < newHeights[stack[len(stack)-1]] {
			curH := newHeights[stack[len(stack)-1]]
			stack = stack[:len(stack)-1]
			curW := i - stack[len(stack)-1] - 1
			ans = max(ans, curH*curW)
		}
		stack = append(stack, i)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
