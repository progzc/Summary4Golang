package leetcode_0088_merge_sorted_array

import "sort"

// 0088.合并两个有序数组
// https://leetcode-cn.com/problems/merge-sorted-array/

// merge_1 使用标准库(即快排)
// 时间复杂度: O((m+n)log(m+n))
// 空间复杂度: O(log(m+n))
func merge_1(nums1 []int, m int, nums2 []int, n int) {
	copy(nums1[m:], nums2)
	sort.Ints(nums1)
}

// merge_2 双指针
// 时间复杂度: O(m+n)
// 空间复杂度: O(m+n)
// 思路: 充分利用讲个数组已经被排序的性质
func merge_2(nums1 []int, m int, nums2 []int, n int) {
	s := make([]int, 0, m+n)
	p1, p2 := 0, 0

	for p1 < m && p2 < n {
		if nums1[p1] <= nums2[p2] {
			s = append(s, nums1[p1])
			p1++
		} else {
			s = append(s, nums2[p2])
			p2++
		}
	}
	if p1 >= m {
		for p2 < n {
			s = append(s, nums2[p2])
			p2++
		}
	}
	if p2 >= n {
		for p1 < m {
			s = append(s, nums1[p1])
			p1++
		}
	}
	copy(nums1, s)
}
