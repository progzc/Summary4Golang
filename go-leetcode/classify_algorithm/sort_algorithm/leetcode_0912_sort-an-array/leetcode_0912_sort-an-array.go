package leetcode_0912_sort_an_array

import (
	"math/rand"
	"time"
)

// 0912. 排序数组
// https://leetcode.cn/problems/sort-an-array/description/

// sortArray 快速排序（超时）
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(h)
func sortArray(nums []int) []int {
	rand.Seed(time.Now().UnixNano())
	quickSort(nums, 0, len(nums)-1)
	return nums
}

func quickSort(nums []int, l, r int) {
	if l >= r {
		return
	}
	q := randomPartition(nums, l, r)
	quickSort(nums, l, q-1)
	quickSort(nums, q+1, r)
}

func randomPartition(a []int, l, r int) int {
	i := rand.Int()%(r-l+1) + l
	a[i], a[r] = a[r], a[i]
	return partition(a, l, r)
}

func partition(a []int, l, r int) int {
	x := a[r]
	i := l - 1
	for j := l; j < r; j++ {
		if a[j] <= x {
			i++
			a[i], a[j] = a[j], a[i]
		}
	}
	a[i+1], a[r] = a[r], a[i+1]
	return i + 1
}
