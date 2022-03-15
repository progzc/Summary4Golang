package leetcode_0260_single_number_iii

// 0260.只出现一次的数字 III
// https://leetcode-cn.com/problems/single-number-iii/
// 概述：有两个元素只出现了一次，其余元素均出现了两次

// singleNumber 位操作
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func singleNumber(nums []int) []int {
	ans := 0
	for _, num := range nums {
		ans ^= num
	}
	div := 1
	for div&ans == 0 {
		div <<= 1
	}
	a, b := 0, 0
	for _, num := range nums {
		if div&num != 0 {
			a ^= num
		} else {
			b ^= num
		}
	}
	return []int{a, b}
}

// singleNumber_2 进一步优化
func singleNumber_2(nums []int) []int {
	ans := 0
	for _, num := range nums {
		ans ^= num
	}
	ans = ans & -ans
	a, b := 0, 0
	for _, num := range nums {
		if ans&num != 0 {
			a ^= num
		} else {
			b ^= num
		}
	}
	return []int{a, b}
}
