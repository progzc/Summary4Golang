package leetcode_0059_spiral_matrix_ii

// 59. 螺旋矩阵 II
// https://leetcode.cn/problems/spiral-matrix-ii/

// generateMatrix 模拟
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：
//	1.首先设定上下左右边界
//	2.其次向右移动到最右，此时第一行因为已经使用过了，可以将其从图中删去，体现在代码中就是重新定义上边界
//	3.判断若重新定义后，上下边界交错，表明螺旋矩阵遍历结束，跳出循环，返回答案
//	4.若上下边界不交错，则遍历还未结束，接着向下向左向上移动，操作过程与第一，二步同理
//	5.不断循环以上步骤，直到某两条边界交错，跳出循环，返回答案
func generateMatrix(n int) [][]int {
	ans := make([][]int, n)
	for i := 0; i < n; i++ {
		ans[i] = make([]int, n)
	}

	up, down, left, right := 0, n-1, 0, n-1
	start := 1
	for true {
		// 左—>右
		for col := left; col <= right; col++ {
			ans[up][col] = start
			start++
		}
		up++
		if up > down {
			break
		}

		// 上—>下
		for row := up; row <= down; row++ {
			ans[row][right] = start
			start++
		}
		right--
		if left > right {
			break
		}

		// 右—>左
		for col := right; col >= left; col-- {
			ans[down][col] = start
			start++
		}
		down--
		if up > down {
			break
		}

		// 下—>上
		for row := down; row >= up; row-- {
			ans[row][left] = start
			start++
		}
		left++
		if left > right {
			break
		}
	}
	return ans
}
