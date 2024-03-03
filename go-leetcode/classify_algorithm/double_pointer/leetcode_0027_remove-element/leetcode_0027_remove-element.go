package leetcode_0027_remove_element

// 0027. 移除元素
// https://leetcode.cn/problems/remove-element/description/

// removeElement 双指针
// 时间复杂度：O(n)
// 空间复杂度：O(1)
// 思路：可以使用双指针：右指针 right 指向当前将要处理的元素，左指针 left 指向下一个将要赋值的位置。
func removeElement(nums []int, val int) int {
	left := 0
	for right := 0; right < len(nums); right++ {
		if nums[right] != val {
			nums[left] = nums[right]
			left++
		}
	}
	return left
}
