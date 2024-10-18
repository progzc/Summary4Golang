package weipai_20241018_water_and_jug_problem

import "fmt"

// 微派四轮面试题
// 0365. 水壶问题
// https://leetcode.cn/problems/water-and-jug-problem

// canMeasureWater 递归
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
func canMeasureWater(x int, y int, target int) bool {
	visit := make(map[string]bool)
	// dfs 含义：表示x壶中剩余 i 升水，y壶中剩余 j 升水，是否可以得到 target 升水
	var dfs func(i, j int) bool
	dfs = func(i, j int) bool {
		key := fmt.Sprintf("%d_%d", i, j)
		if visit[key] {
			return false
		}
		visit[key] = true
		if i == target || j == target || i+j == target {
			return true
		}
		// 主要分 6 种情况：
		// 1. dfs(x,j)表示把 x 壶灌满
		// 2. dfs(i,y)表示把 y 壶灌满
		// 3. dfs(0,j)表示把 x 壶倒空
		// 4. dfs(i,0)表示把 y 壶倒空
		if dfs(x, j) || dfs(i, y) || dfs(0, j) || dfs(i, 0) {
			return true
		}
		a := min(i, y-j)
		b := min(j, x-i)
		// 5. dfs(i-a, j+a)表示把 x 壶的水灌满 y 壶，直至灌满或倒空
		// 6. dfs(i+b, j-b)表示把 y 壶的水灌满 x 壶，直至倒满或倒空
		return dfs(i-a, j+a) || dfs(i+b, j-b)
	}
	return dfs(0, 0)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
