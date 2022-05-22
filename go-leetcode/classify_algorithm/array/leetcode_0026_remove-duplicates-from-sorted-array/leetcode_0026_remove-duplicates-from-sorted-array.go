package leetcode_0026_remove_duplicates_from_sorted_array

// 0026.删除有序数组中的重复项
// https://leetcode.cn/problems/remove-duplicates-from-sorted-array/

// removeDuplicates
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func removeDuplicates(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return n
	}

	j := 0
	for i := 1; i < n; i++ {
		if nums[i-1] == nums[i] {
			continue
		}
		j++
		nums[j] = nums[i]
	}
	return j + 1
}

// removeDuplicates_2 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func removeDuplicates_2(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return n
	}

	slow := 1
	for fast := 1; fast < n; fast++ {
		if nums[fast-1] != nums[fast] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
