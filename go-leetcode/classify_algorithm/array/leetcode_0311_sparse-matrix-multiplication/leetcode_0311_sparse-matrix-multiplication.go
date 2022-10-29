package leetcode_0311_sparse_matrix_multiplication

// 0311. 稀疏矩阵的乘法
// https://leetcode.cn/problems/sparse-matrix-multiplication/

// multiply 直接计算
// 时间复杂度: O(m*n*k)
// 空间复杂度: O(mn)
func multiply(mat1 [][]int, mat2 [][]int) [][]int {
	var (
		m, k, n = len(mat1), len(mat1[0]), len(mat2[0])
		ans     = make([][]int, m)
	)
	for i := 0; i < m; i++ {
		ans[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			for p := 0; p < k; p++ {
				ans[i][j] += mat1[i][p] * mat2[p][j]
			}
		}
	}
	return ans
}

// multiply_2 直接计算（优化）
// 时间复杂度: O(m*n*k)
// 空间复杂度: O(mn)
func multiply_2(mat1 [][]int, mat2 [][]int) [][]int {
	var (
		m, k, n = len(mat1), len(mat1[0]), len(mat2[0])
		ans     = make([][]int, m)
	)
	for i := 0; i < m; i++ {
		ans[i] = make([]int, n)
	}

	for i := 0; i < m; i++ {
		for p := 0; p < k; p++ {
			if mat1[i][p] == 0 {
				continue
			}
			for j := 0; j < n; j++ {
				ans[i][j] += mat1[i][p] * mat2[p][j]
			}
		}
	}
	return ans
}

// multiply_3 工程优化（空间换时间）
// 时间复杂度: O(m*n)
// 空间复杂度: O(mn)
func multiply_3(mat1 [][]int, mat2 [][]int) [][]int {
	var (
		m, n = len(mat1), len(mat2[0])
		ans  = make([][]int, m)
	)
	for i := 0; i < m; i++ {
		ans[i] = make([]int, n)
	}

	noneZeroA := getNoneZeroMat(mat1)
	noneZeroB := getNoneZeroMat(mat2)
	for _, m1 := range noneZeroA {
		for _, m2 := range noneZeroB {
			// 这里这么判断的原因：mat1和mat2相乘，
			// 只有mat1的数据的列数和mat2的行数相等才会进行计算
			if m1[1] == m2[0] {
				ans[m1[0]][m2[1]] += m1[2] * m2[2]
			}
		}
	}
	return ans
}

func getNoneZeroMat(matrix [][]int) [][]int {
	m, n := len(matrix), len(matrix[0])
	var ans [][]int
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] != 0 {
				ans = append(ans, []int{i, j, matrix[i][j]})
			}
		}
	}
	return ans
}
