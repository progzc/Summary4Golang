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
// 思路: 充分利用两个数组已经被排序的性质
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

// merge_3 逆向双指针
// 时间复杂度: O(m+n)
// 空间复杂度: O(1)
// 思路: 充分利用两个数组已经被排序的性质 + nums1的长度为m+n
func merge_3(nums1 []int, m int, nums2 []int, n int) {
	p1, p2, i := m-1, n-1, m+n-1

	for p1 >= 0 && p2 >= 0 {
		if nums1[p1] <= nums2[p2] {
			nums1[i] = nums2[p2]
			p2--
		} else {
			nums1[i] = nums1[p1]
			p1--
		}
		i--
	}
	if p1 < 0 {
		for p2 >= 0 {
			nums1[i] = nums2[p2]
			p2--
			i--
		}
	}
	if p2 < 0 {
		for p1 >= 0 {
			nums1[i] = nums1[p1]
			p1--
			i--
		}
	}
}
