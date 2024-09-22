package leetcode_0283_move_zeroes

// 0283.移动零
// https://leetcode.cn/problems/move-zeroes/description/

// moveZeroes 双指针
// 思路：
// 使用双指针，左指针指向当前已经处理好的序列的尾部，右指针指向待处理序列的头部。
// 右指针不断向右移动，每次右指针指向非零数，则将左右指针对应的数交换，同时左指针右移。
func moveZeroes(nums []int) {
	left, right, n := 0, 0, len(nums)
	for right < n {
		if nums[right] != 0 {
			nums[left], nums[right] = nums[right], nums[left]
			left++
			right++
		} else {
			right++
		}
	}
}
