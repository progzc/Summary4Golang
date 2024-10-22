package leetcode_0905_sort_array_by_parity

// 0905. 按奇偶排序数组
// https://leetcode.cn/problems/sort-array-by-parity/

// sortArrayByParity
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func sortArrayByParity(nums []int) []int {
	n := len(nums)
	if n <= 1 {
		return nums
	}
	p1, p2 := 0, n-1
	for p1 < p2 {
		for p1 < p2 && nums[p1]%2 == 0 {
			p1++
		}
		for p1 < p2 && nums[p2]%2 == 1 {
			p2--
		}
		if p1 < p2 {
			nums[p1], nums[p2] = nums[p2], nums[p1]
			p1++
			p2--
		}
	}
	return nums
}
