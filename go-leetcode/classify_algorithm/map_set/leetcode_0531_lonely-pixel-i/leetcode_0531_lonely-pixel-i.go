package leetcode_0531_lonely_pixel_i

// 0531. 孤独像素 I
// https://leetcode.cn/problems/lonely-pixel-i/

// findLonelyPixel 暴力法
// 时间复杂度: O(mn(m+n))
// 空间复杂度: O(m+n)
func findLonelyPixel(picture [][]byte) int {
	if len(picture) == 0 || len(picture[0]) == 0 {
		return 0
	}

	m, n := len(picture), len(picture[0])
	ans := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if picture[i][j] == 'B' {
				exist := true
				// 检查横向
				for k := 0; k < n && exist; k++ {
					if k != j && picture[i][k] == 'B' {
						exist = false
					}
				}
				// 检查纵向
				for k := 0; k < m && exist; k++ {
					if k != i && picture[k][j] == 'B' {
						exist = false
					}
				}
				if exist {
					ans++
				}
			}
		}
	}
	return ans
}

// findLonelyPixel_2 哈希
// 时间复杂度: O(mn)
// 空间复杂度: O(m+n)
func findLonelyPixel_2(picture [][]byte) int {
	x, y := map[int]int{}, map[int]int{}
	for i := 0; i < len(picture); i++ {
		for j := 0; j < len(picture[0]); j++ {
			if picture[i][j] == 'B' {
				x[i]++
				y[j]++
			}
		}
	}

	ans := 0
	for i := 0; i < len(picture); i++ {
		for j := 0; j < len(picture[0]); j++ {
			if picture[i][j] == 'B' && x[i] == 1 && y[j] == 1 {
				ans++
			}
		}
	}
	return ans
}
