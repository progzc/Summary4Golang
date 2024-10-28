package leetcode_0075_sort_colors

import (
	"fmt"
	"testing"
)

// 0075. 颜色分类🌟
// https://leetcode.cn/problems/sort-colors

func TestSortColors(t *testing.T) {
	nums := []int{1, 1, 2, 0, 2, 2}
	sortColors(nums)
	fmt.Println(nums) // [0 1 1 2 2 2]
}

// sortColors 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func sortColors(nums []int) {
	p0, p2 := 0, len(nums)-1
	for i := 0; i <= p2; i++ {
		for i <= p2 && nums[i] == 2 {
			nums[i], nums[p2] = nums[p2], nums[i]
			p2--
		}
		if nums[i] == 0 {
			nums[i], nums[p0] = nums[p0], nums[i]
			p0++
		}
	}
}
