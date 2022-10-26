package leetcode_1182_shortest_distance_to_target_color

import (
	"math"
	"sort"
)

// 1182. 与目标颜色间的最短距离
// https://leetcode.cn/problems/shortest-distance-to-target-color/

// shortestDistanceColor 二分
// 时间复杂度: O(q*log(n))
// 空间复杂度: O(n)
func shortestDistanceColor(colors []int, queries [][]int) []int {
	m := make(map[int][]int)
	for i, c := range colors {
		m[c] = append(m[c], i)
	}
	var ans []int
	for _, q := range queries {
		dis := m[q[1]]
		if len(dis) == 0 {
			ans = append(ans, -1)
			continue
		}
		// 二分查找找到目标数字插入的位置
		d := math.MaxInt32
		idx := sort.SearchInts(dis, q[0])
		/*
			if idx == 0 {
				d = min(d, dis[idx]-q[0])
			} else if idx == len(dis) {
				d = min(d, q[0]-dis[idx-1])
			} else {
				d = min(min(d, dis[idx]-q[0]), q[0]-dis[idx-1])
			}
		*/
		// 上面这段注释的代码可以简写为
		if idx < len(dis) {
			d = min(d, dis[idx]-q[0])
		}
		if idx-1 >= 0 {
			d = min(d, q[0]-dis[idx-1])
		}

		ans = append(ans, d)
	}
	return ans
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
