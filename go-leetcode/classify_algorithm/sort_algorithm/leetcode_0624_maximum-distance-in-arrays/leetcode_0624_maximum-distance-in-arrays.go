package leetcode_0624_maximum_distance_in_arrays

// 0624. 数组列表中的最大距离
// https://leetcode.cn/problems/maximum-distance-in-arrays/

// maxDistance 暴力法(超时)
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
func maxDistance(arrays [][]int) int {
	ans := 0
	n := len(arrays)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			l1, l2 := len(arrays[i]), len(arrays[j])
			ans = max(ans, abs(arrays[i][0]-arrays[j][l2-1]))
			ans = max(ans, abs(arrays[i][l1-1]-arrays[j][0]))
		}
	}
	return ans
}

// maxDistance_2 线性扫描
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func maxDistance_2(arrays [][]int) int {
	ans := 0
	n := len(arrays)
	minVal, maxVal := arrays[0][0], arrays[0][len(arrays[0])-1]
	for i := 1; i < n; i++ {
		ans = max(ans, abs(arrays[i][len(arrays[i])-1]-minVal))
		ans = max(ans, abs(arrays[i][0]-maxVal))
		minVal = min(minVal, arrays[i][0])
		maxVal = max(maxVal, arrays[i][len(arrays[i])-1])
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
