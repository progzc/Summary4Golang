package leetcode_0562_longest_line_of_consecutive_one_in_matrix

// 0562. 矩阵中最长的连续1线段
// https://leetcode.cn/problems/longest-line-of-consecutive-one-in-matrix/

// longestLine 暴力法(通过剪枝速度仍然很快)
// 时间复杂度: O((mn)^2)
// 空间复杂度: O(1)
func longestLine(mat [][]int) int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return 0
	}

	row, column := len(mat), len(mat[0])
	ans := 0
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			if mat[i][j] == 0 {
				continue
			}

			// 横向检查（剪枝）
			if column-j > ans {
				count := 1
				for m := j + 1; m < column; m++ {
					if mat[i][m] == 0 {
						break
					}
					count++
				}
				ans = max(ans, count)
			}

			// 纵向检查（剪枝）
			if row-i > ans {
				count := 1
				for n := i + 1; n < row; n++ {
					if mat[n][j] == 0 {
						break
					}
					count++
				}
				ans = max(ans, count)
			}

			// 对角线检查（剪枝）
			if row-i > ans && column-j > ans {
				count := 1
				for k := 1; i+k < row && j+k < column; k++ {
					if mat[i+k][j+k] == 0 {
						break
					}
					count++
				}
				ans = max(ans, count)
			}

			// 反对角线检查（剪枝）
			if row-i > ans && j+1 > ans {
				count := 1
				for k := 1; i+k < row && j-k >= 0; k++ {
					if mat[i+k][j-k] == 0 {
						break
					}
					count++
				}
				ans = max(ans, count)
			}
		}
	}
	return ans
}

// longestLine_2 二维动态规划
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路：
//	设 dp[i][j] 代表以矩阵第 i 行，第 j 列结尾的连续 1 线段的长度。
//	那么我们需要记录四个 dp 矩阵：横，竖，对角线，以及反对角线分别需要一个 dp 矩阵来记录。
//		横向 dp 矩阵仅需要考虑其左边的连续 1 线段，即 dp[i][j - 1]。
//		纵向 dp 矩阵仅需要考虑其上边的连续 1 线段，即 dp[i - 1][j]。
//		对角线方向 dp 矩阵仅需要考虑其左上的连续 1 线段，即 dp[i - 1][j - 1]。
//		反对角线方向 dp 矩阵仅需要考虑其右上的连续 1 线段，即 dp[i - 1][j + 1]。
func longestLine_2(mat [][]int) int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return 0
	}

	row, column := len(mat), len(mat[0])
	// 水平/垂直/对角/反对角
	h, v, d, a := make([][]int, row), make([][]int, row), make([][]int, row), make([][]int, row)
	for i := 0; i < row; i++ {
		h[i], v[i], d[i], a[i] = make([]int, column), make([]int, column), make([]int, column), make([]int, column)
	}

	ans := 0
	for i := 0; i < row; i++ {
		for j := 0; j < column; j++ {
			if mat[i][j] == 0 {
				h[i][j], v[i][j], d[i][j], a[i][j] = 0, 0, 0, 0
			} else {
				// 水平方向
				if j > 0 {
					h[i][j] = h[i][j-1] + 1
				} else {
					h[i][j] = 1
				}
				// 垂直方向
				if i > 0 {
					v[i][j] = v[i-1][j] + 1
				} else {
					v[i][j] = 1
				}
				// 对角方向
				if i > 0 && j > 0 {
					d[i][j] = d[i-1][j-1] + 1
				} else {
					d[i][j] = 1
				}
				// 反对角方向
				if i > 0 && j < column-1 {
					a[i][j] = a[i-1][j+1] + 1
				} else {
					a[i][j] = 1
				}
				ans = max(ans, h[i][j])
				ans = max(ans, v[i][j])
				ans = max(ans, d[i][j])
				ans = max(ans, a[i][j])
			}
		}
	}
	return ans
}

// longestLine_3 一维动态规划
// 时间复杂度: O(mn)
// 空间复杂度: O(n)
// 思路：
//	设 dp[i][j] 代表以矩阵第 i 行，第 j 列结尾的连续 1 线段的长度。
//	那么我们需要记录四个 dp 矩阵：横，竖，对角线，以及反对角线分别需要一个 dp 矩阵来记录。
//		横向 dp 矩阵仅需要考虑其左边的连续 1 线段，即 dp[i][j - 1]。
//		纵向 dp 矩阵仅需要考虑其上边的连续 1 线段，即 dp[i - 1][j]。
//		对角线方向 dp 矩阵仅需要考虑其左上的连续 1 线段，即 dp[i - 1][j - 1]。
//		反对角线方向 dp 矩阵仅需要考虑其右上的连续 1 线段，即 dp[i - 1][j + 1]。
// 注意：
//	dp 数组中的每一个位置的值只依赖于上一行。 因此不需要将整个矩阵的结果全部存储，只需要保留上一行的结果即可。
// 	而对于横向的 dp 数组，由于其不依赖于上一行，上一行的结果也可以不存储。这样可以达到节省空间的效果。
func longestLine_3(mat [][]int) int {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return 0
	}

	row, column := len(mat), len(mat[0])
	// 水平/垂直/对角/反对角
	h, v, d, a := make([]int, column), make([]int, column), make([]int, column), make([]int, column)

	ans := 0
	for i := 0; i < row; i++ {
		vn, dn, an := make([]int, column), make([]int, column), make([]int, column)
		for j := 0; j < column; j++ {
			if mat[i][j] == 0 {
				h[j], vn[j], dn[j], an[j] = 0, 0, 0, 0
			} else {
				// 水平方向
				if j > 0 {
					h[j] = h[j-1] + 1
				} else {
					h[j] = 1
				}
				// 垂直方向
				if i > 0 {
					vn[j] = v[j] + 1
				} else {
					vn[j] = 1
				}
				// 对角方向
				if i > 0 && j > 0 {
					dn[j] = d[j-1] + 1
				} else {
					dn[j] = 1
				}
				// 反对角方向
				if i > 0 && j < column-1 {
					an[j] = a[j+1] + 1
				} else {
					an[j] = 1
				}
				ans = max(ans, h[j])
				ans = max(ans, vn[j])
				ans = max(ans, dn[j])
				ans = max(ans, an[j])
			}
		}
		v, d, a = vn, dn, an
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
