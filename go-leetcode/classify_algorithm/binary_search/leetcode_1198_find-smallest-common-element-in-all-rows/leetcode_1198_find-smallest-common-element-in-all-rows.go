package leetcode_1198_find_smallest_common_element_in_all_rows

import "sort"

// 1198. 找出所有行中最小公共元素
// https://leetcode.cn/problems/find-smallest-common-element-in-all-rows/

// smallestCommonElement 二分
// 时间复杂度: O(n*m*log(m))
// 空间复杂度: O(1)
// 思路：
//	a.遍历第一行所有元素，使用二分搜索判断其他行是否都存在该元素。
//	b.改进的方法：如果从上一次搜索返回位置开始搜索可以降低平均时间复杂度。如果一行所有元素都小于查找值，则返回 -1。
func smallestCommonElement(mat [][]int) int {
	m, n := len(mat), len(mat[0])
	pos := make([]int, m)
	for j := 0; j < n; j++ {
		find := true
		for i := 1; i < m; i++ {
			pos[i] = binarySearch(mat[i], mat[0][j])
			if pos[i] < 0 {
				find = false
				break
			}
		}
		if find {
			return mat[0][j]
		}
	}
	return -1
}

func binarySearch(nums []int, target int) int {
	low, high := 0, len(nums)-1
	for low <= high {
		mid := low + (high-low)>>1
		if nums[mid] == target {
			return mid
		} else if nums[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

// smallestCommonElement_2 哈希计数
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路：
//	a.统计元素出现的次数。
//	b.逐列计算元素，可以提高平均时间复杂度。这样，首先计算较小元素的出现次数，一旦某个元素出现次数为 n，则直接返回。
func smallestCommonElement_2(mat [][]int) int {
	count := make([]int, 10001)
	m, n := len(mat), len(mat[0])
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			count[mat[i][j]]++
			if count[mat[i][j]] == m {
				return mat[i][j]
			}
		}
	}
	return -1
}

// smallestCommonElement_3 哈希计数
// 时间复杂度: O(mn)
// 空间复杂度: O(mn)
// 思路：
//	a.统计元素出现的次数。
func smallestCommonElement_3(mat [][]int) int {
	count := make(map[int]int)
	vals := make([]int, 0)
	m, n := len(mat), len(mat[0])
	for j := 0; j < n; j++ {
		for i := 0; i < m; i++ {
			if v, ok := count[mat[i][j]]; ok {
				count[mat[i][j]] = v + 1
			} else {
				count[mat[i][j]] = 1
				vals = append(vals, mat[i][j])
			}
		}
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})

	for _, val := range vals {
		if count[val] == m {
			return val
		}
	}
	return -1
}
