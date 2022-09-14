package leetcode_0469_convex_polygon

// 0469. 凸多边形
// https://leetcode.cn/problems/convex-polygon/

// 数学知识：https://blog.csdn.net/houyichaochao/article/details/81141893

// isConvex 叉乘
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func isConvex(points [][]int) bool {
	n := len(points)
	pre := 0
	for i := 0; i < n; i++ {
		x1 := points[(i+1)%n][0] - points[i][0]
		y1 := points[(i+1)%n][1] - points[i][1]

		x2 := points[(i+2)%n][0] - points[(i+1)%n][0]
		y2 := points[(i+2)%n][1] - points[(i+1)%n][1]

		// 叉乘
		if multi := x1*y2 - x2*y1; multi != 0 {
			// 与上一次的法向量方向相反
			if pre*multi < 0 {
				return false
			}
			pre = multi
		}
	}
	return true
}
