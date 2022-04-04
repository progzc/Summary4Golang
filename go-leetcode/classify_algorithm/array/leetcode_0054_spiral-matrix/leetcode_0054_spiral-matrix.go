package leetcode_0054_spiral_matrix

// 0054.螺旋矩阵
// https://leetcode-cn.com/problems/spiral-matrix/

// spiralOrder 模拟
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
// 思路：方向向量+终止条件
func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return []int{}
	}

	rows, columns := len(matrix), len(matrix[0])
	total := rows * columns
	m := make([][]bool, rows)
	for i := 0; i < rows; i++ {
		m[i] = make([]bool, columns)
	}

	row, column := 0, 0
	di := 0
	dv := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	ans := make([]int, rows*columns)
	// 终止条件
	for i := 0; i < total; i++ {
		ans[i] = matrix[row][column]
		m[row][column] = true
		nextRow, nextColumn := row+dv[di][0], column+dv[di][1]
		// 决定是否调转方向
		if nextRow < 0 || nextRow >= rows || nextColumn < 0 || nextColumn >= columns || m[nextRow][nextColumn] {
			// 方向矢量
			di = (di + 1) % 4
		}
		row, column = row+dv[di][0], column+dv[di][1]
	}
	return ans
}

// spiralOrder_2 按圈模拟
// 时间复杂度: O(m*n)
// 空间复杂度: O(1)
// 思路：计算推演下一个上下左右的位置
//	左上：top,left
//	右上：top,right
//	右下：bottom,right
//	左下：bottom,left
//	需要注意最后一圈的特殊情况
func spiralOrder_2(matrix [][]int) []int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return []int{}
	}

	rows, columns := len(matrix), len(matrix[0])
	ans := make([]int, rows*columns)
	left, right, top, bottom := 0, columns-1, 0, rows-1
	i := 0

	for left <= right && top <= bottom {
		// 从左到右
		for column := left; column <= right; column++ {
			ans[i] = matrix[top][column]
			i++
		}
		// 从上到下
		for row := top + 1; row <= bottom; row++ {
			ans[i] = matrix[row][right]
			i++
		}

		// 注意：考虑奇数行列最后一圈的走法（或者单独考虑一行、一列的特例）
		// 每次都是先往右走和先往下走，所以往右走和往下走的路一定是没走过的路，不需要加if判断；
		// 往左走的时候，如果top == bottom，那么会重复走之前从左往右走过的路，所以需要加上top < bottom的判断，同理往上走也一样。
		// 不然的话，重复走会导致ans[i]中的i索引越界
		if left < right && top < bottom {
			// 从右到左
			for column := right - 1; column >= left; column-- {
				ans[i] = matrix[bottom][column]
				i++
			}
			// 从下到上
			for row := bottom - 1; row >= top+1; row-- {
				ans[i] = matrix[row][left]
				i++
			}
		}
		// 进行下一圈
		left++
		right--
		top++
		bottom--
	}
	return ans
}
