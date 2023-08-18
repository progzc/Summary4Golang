package leetcode_0280_wiggle_sort

import "sort"

// 0280. 摆动排序
// https://leetcode.cn/problems/wiggle-sort/

// wiggleSort 排序
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(n)
func wiggleSort(nums []int) {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	n := len(nums)
	l, r := 0, n-1
	res := make([]int, n)
	flag := true
	for i := 0; i < n; i++ {
		if flag {
			res[i] = nums[l]
			l++
		} else {
			res[i] = nums[r]
			r--
		}
		flag = !flag
	}

	for i := 0; i < n; i++ {
		nums[i] = res[i]
	}
}

// wiggleSort_2 排序
// 时间复杂度: O(n*log(n))
// 空间复杂度: O(1)
// 思路: 先将数组排序，再从第二个元素开始逐对交换相邻元素的位置
func wiggleSort_2(nums []int) {
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	for i := 1; i < len(nums)-1; i += 2 {
		nums[i], nums[i+1] = nums[i+1], nums[i]
	}
}

// wiggleSort_3 一次遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 数学归纳法证明:
//	假设 [0,1, ..., k] 均已满足摆动排序，而 k+1 不满足，记此时 nums[k-1]，nums[k]，nums[k+1] 分别为 a，b，c
// 	若不满足的是降序，则可知：c>b，同时 b>=a，此时三个数状态是 a <= b < c，交换 b 和 c 后变成：a < c > b，满足小大小
//	若不满足的是升序，则可知：c<b，同时 b<=a，此时三个数状态是 a >= b > c，交换 b 和 c 后变成：a > c < b，满足大小大
func wiggleSort_3(nums []int) {
	less := true
	for i := 0; i < len(nums)-1; i++ {
		if less {
			if nums[i] > nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
			}
		} else {
			if nums[i] < nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
			}
		}
		less = !less
	}
}

// wiggleSort_4 摆动排序
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func wiggleSort_4(nums []int)  {
	for i:=1; i<len(nums); i=i+2 {
		if nums[i-1]>nums[i] {
			nums[i-1],nums[i] = nums[i],nums[i-1]
		}
		if i+1<len(nums) && nums[i]<nums[i+1] {
			nums[i+1],nums[i] = nums[i],nums[i+1]
		}
	}
}
