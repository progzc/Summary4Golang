package leetcode_0750_number_of_corner_rectangles

// 0750. 角矩形的数量
// https://leetcode.cn/problems/number-of-corner-rectangles/

// countCornerRectangles 常规方法
// 时间复杂度: O((mn)^2)
// 空间复杂度: O(1)
func countCornerRectangles(grid [][]int) int {
	if grid == nil || len(grid) == 0 {
		return 0
	}

	row, column := len(grid), len(grid[0])
	isLands := 0

	for i := 0; i < row-1; i++ {
		for j := 0; j < column-1; j++ {
			if grid[i][j] == 0 {
				continue
			}
			for m := 1; m < row-i; m++ {
				if grid[i+m][j] == 0 {
					continue
				}
				for n := 1; n < column-j; n++ {
					if grid[i][j+n] == 1 && grid[i+m][j+n] == 1 {
						isLands++
					}
				}
			}
		}
	}
	return isLands
}
