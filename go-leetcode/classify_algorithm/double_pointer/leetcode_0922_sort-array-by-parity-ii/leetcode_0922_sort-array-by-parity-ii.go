package leetcode_0922_sort_array_by_parity_ii

// 0922. 按奇偶排序数组 II
// https://leetcode.cn/problems/sort-array-by-parity-ii

// sortArrayByParityII
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func sortArrayByParityII(nums []int) []int {
	n := len(nums)
	if n%2 == 1 || n <= 1 {
		return nums
	}
	p1, p2 := 0, 1
	for p1 < n && p2 < n {
		for p1 < n && nums[p1]%2 == 0 {
			p1 += 2
		}
		for p2 < n && nums[p2]%2 == 1 {
			p2 += 2
		}
		if p1 < n && p2 < n {
			nums[p1], nums[p2] = nums[p2], nums[p1]
		}
	}
	return nums
}
