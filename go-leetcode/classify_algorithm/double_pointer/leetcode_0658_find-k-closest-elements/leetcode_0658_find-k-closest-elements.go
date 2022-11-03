package leetcode_0658_find_k_closest_elements

import (
	"sort"
)

// 0658. 找到 K 个最接近的元素
// https://leetcode.cn/problems/find-k-closest-elements/

// findClosestElements 排序法
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(log(n))
// 注意: 并不要求连续
// 思路:
//	首先将数组 arr 按照「更接近」的定义进行排序，如果 a 比 b 更接近 x，那么 a 将排在 b 前面。
//	排序完成之后，k 个最接近的元素就是数组 arr 的前 k 个元素，将这 k 个元素从小到大进行排序后，直接返回。
func findClosestElements(arr []int, k int, x int) []int {
	sort.SliceStable(arr, func(i, j int) bool {
		return abs(arr[i]-x) < abs(arr[j]-x)
	})
	arr = arr[:k]
	sort.Ints(arr)
	return arr
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

// findClosestElements_2 二分+双指针
// 时间复杂度: O(log(n)+k)
// 空间复杂度: O(1)
// 注意: 并不要求连续
// 排序法的缺点: 并没有利用arr数组升序排列这个特点
func findClosestElements_2(arr []int, k int, x int) []int {
	n := len(arr)
	// 关键是理解下面这两行
	r := sort.SearchInts(arr, x)
	l := r - 1
	for k > 0 {
		if l < 0 {
			r++
		} else if r >= n || x-arr[l] <= arr[r]-x {
			l--
		} else {
			r++
		}
		k--
	}
	return arr[l+1 : r]
}
