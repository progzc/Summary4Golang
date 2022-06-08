package leetcode_0163_missing_ranges

import (
	"fmt"
	"strings"
)

// 0163.缺失的区间
// https://leetcode.cn/problems/missing-ranges/

// findMissingRanges
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func findMissingRanges(nums []int, lower int, upper int) []string {
	var (
		ans []string
		n   = len(nums)
	)
	if n == 0 {
		ans = append(ans, helper(lower-1, upper+1))
		return ans
	}

	if lower < nums[0] {
		ans = append(ans, helper(lower-1, nums[0]))
	}
	for i := 0; i < n-1; i++ {
		if nums[i]+1 != nums[i+1] {
			ans = append(ans, helper(nums[i], nums[i+1]))
		}
	}
	if nums[n-1] < upper {
		ans = append(ans, helper(nums[n-1], upper+1))
	}
	return ans
}

func helper(left, right int) string {
	sb := strings.Builder{}
	if left+2 == right {
		sb.WriteString(fmt.Sprintf("%d", left+1))
	} else {
		sb.WriteString(fmt.Sprintf("%d->%d", left+1, right-1))
	}
	return sb.String()
}

// findMissingRanges_2
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func findMissingRanges_2(nums []int, lower int, upper int) []string {
	var (
		ans []string
	)
	// 预处理，这样就不需要考虑边界问题了
	nums = append([]int{lower - 1}, nums...)
	nums = append(nums, upper+1)

	n := len(nums)
	for i := 0; i < n-1; i++ {
		if nums[i]+1 != nums[i+1] {
			left, right := nums[i]+1, nums[i+1]-1
			if left == right {
				ans = append(ans, fmt.Sprintf("%d", left))
			} else {
				ans = append(ans, fmt.Sprintf("%d->%d", left, right))
			}
		}
	}
	return ans
}
