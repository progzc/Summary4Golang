package leetcode_0015_3sum

import (
	"sort"
)

// 0015.三数之和
// https://leetcode-cn.com/problems/3sum/

// threeSum 排序 + 双指针
// 时间复杂度: O(n^2)
// 空间复杂度: O(log(n))
// 思路：
// 1. 凡是不重复，首先要想到先排序 或 采用散列表
// 2. 排序之后采用三重循环，要保证每重循环 相邻两次枚举的元素不能相同
// 双指针：当随着第一个元素的递增,第二个元素是递减的，就可以考虑使用 双指针
func threeSum(nums []int) [][]int {
	n := len(nums)
	sort.Ints(nums)
	ans := [][]int{}

	// 枚举a
	for first := 0; first < n; first++ {
		// 需要和上一次枚举的数不同(但是第一次循环必须执行)
		if first > 0 && nums[first] == nums[first-1] {
			continue
		}
		third := n - 1
		target := -1 * nums[first]
		// 枚举b
		for second := first + 1; second < n; second++ {
			// 需要和上一次枚举的数不同(但是第一次循环必须执行)
			if second > first+1 && nums[second] == nums[second-1] {
				continue
			}
			// 保证b在c左侧
			for second < third && nums[second]+nums[third] > target {
				third--
			}
			// 提前退出
			if second == third {
				break
			}
			if nums[second]+nums[third] == target {
				ans = append(ans, []int{nums[first], nums[second], nums[third]})
			}
		}
	}
	return ans
}
