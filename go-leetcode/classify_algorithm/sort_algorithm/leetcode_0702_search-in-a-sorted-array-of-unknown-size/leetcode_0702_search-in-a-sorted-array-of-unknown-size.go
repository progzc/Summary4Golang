package leetcode_0702_search_in_a_sorted_array_of_unknown_size

import "math"

// 0702. 搜索长度未知的有序数组
// https://leetcode.cn/problems/search-in-a-sorted-array-of-unknown-size/

type ArrayReader struct {
	nums []int
}

func (this *ArrayReader) get(index int) int {
	if index > len(this.nums)-1 {
		return math.MaxInt32
	}
	return this.nums[index]
}

func search(reader ArrayReader, target int) int {
	if reader.get(0) == target {
		return 0
	}

	// 搜索右边界
	left, right := 0, 1
	for reader.get(right) < target {
		left = right
		right <<= 1
	}

	// 二分搜索
	var (
		pivot int
		num   int
	)
	for left <= right {
		pivot = left + (right-left)>>1
		num = reader.get(pivot)
		if num == target {
			return pivot
		} else if num > target {
			right = pivot - 1
		} else {
			left = pivot + 1
		}
	}
	return -1
}
