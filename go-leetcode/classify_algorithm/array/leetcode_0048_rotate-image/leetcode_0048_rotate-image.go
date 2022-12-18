package leetcode_0048_rotate_image

// 48. 旋转图像
// https://leetcode.cn/problems/rotate-image/

// rotate 使用辅助数组
// 时间复杂度: O(n^2)
// 空间复杂度: O(n^2)
// 思路:
//	观察可知, 旋转后的规律如下:
//		对于矩阵中第 i 行的第 j 个元素，在旋转后，它出现在倒数第 i 列的第 j 个位置。
//		matrix[row][col] --> matrix[col][n-row-1]
func rotate(matrix [][]int) {
	n := len(matrix)
	newMatrix := make([][]int, n)
	for i := 0; i < n; i++ {
		newMatrix[i] = make([]int, n)
	}
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			newMatrix[col][n-row-1] = matrix[row][col]
		}
	}
	copy(matrix, newMatrix)
}

// rotate_2 原地旋转
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
// 思路:
//	观察可知, 旋转后的规律如下:
//		对于矩阵中第 i 行的第 j 个元素，在旋转后，它出现在倒数第 i 列的第 j 个位置。
//		matrix[row][col] --> matrix[col][n-row-1]
//	反复应用上面的公式即可求解:
//		1.第一步: matrix[row][col] --> matrix[col][n-row-1]
//		2.第二步: matrix[col][n-row-1] --> matrix[n-row-1][n-col-1]
//		3.第三步: matrix[n-row-1][n-col-1] --> matrix[n-col-1][row]
//		4.第四步: matrix[n-col-1][row] --> matrix[row][col]
//		上述步骤的逆过程即为旋转交换的顺序。
func rotate_2(matrix [][]int) {
	n := len(matrix)
	for row := 0; row < n/2; row++ {
		for col := 0; col < (n+1)/2; col++ {
			matrix[row][col], matrix[n-col-1][row], matrix[n-row-1][n-col-1], matrix[col][n-row-1] =
				matrix[n-col-1][row], matrix[n-row-1][n-col-1], matrix[col][n-row-1], matrix[row][col]
		}
	}
}

// rotate_3 用翻转代替旋转
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
// 思路:
//	观察可知, 旋转后的规律如下:
//		对于矩阵中第 i 行的第 j 个元素，在旋转后，它出现在倒数第 i 列的第 j 个位置。
//		matrix[row][col] --> matrix[col][n-row-1]
//	可以使用翻转代替旋转:
//		第一步: 上下翻转, matrix[row][col] --> matrix[n-row-1][col]
//		第二步: 主对角线翻转, matrix[n-row-1][col] --> matrix[col][n-row-1]
func rotate_3(matrix [][]int) {
	n := len(matrix)
	// 上下翻转
	for row := 0; row < n/2; row++ {
		matrix[row], matrix[n-1-row] = matrix[n-1-row], matrix[row]
	}
	// 主对角线翻转
	for row := 0; row < n; row++ {
		for col := 0; col < row; col++ {
			matrix[row][col], matrix[col][row] = matrix[col][row], matrix[row][col]
		}
	}
}
