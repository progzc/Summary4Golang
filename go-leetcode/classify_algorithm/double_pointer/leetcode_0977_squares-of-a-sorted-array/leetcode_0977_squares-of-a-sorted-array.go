package leetcode_0977_squares_of_a_sorted_array

import "sort"

// 0977. 有序数组的平方
// https://leetcode.cn/problems/squares-of-a-sorted-array/description/

// sortedSquares 双指针
// 时间复杂度：O(n)
// 空间复杂度: O(1)
func sortedSquares(nums []int) []int {
	ans := make([]int, len(nums))
	left, right := 0, len(nums)-1
	i := len(nums) - 1
	for left <= right {
		if nums[left]*nums[left] >= nums[right]*nums[right] {
			ans[i] = nums[left] * nums[left]
			left++
		} else {
			ans[i] = nums[right] * nums[right]
			right--
		}
		i--
	}
	return ans
}

// sortedSquares_2 排序（不会超时）
// 时间复杂度：O(nlog(n))
// 空间复杂度: O(log(n))
func sortedSquares_2(nums []int) []int {
	ans := make([]int, len(nums))
	for i, num := range nums {
		ans[i] = num * num
	}
	sort.Ints(ans)
	return ans
}
