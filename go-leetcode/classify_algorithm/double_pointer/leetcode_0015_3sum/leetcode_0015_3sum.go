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
	var ans [][]int
	n := len(nums)
	if n < 3 {
		return ans
	}

	sort.Ints(nums)
	// 枚举a
	for i := 0; i < n-2; i++ {
		// 需要和上一次枚举的数不同(但是第一次循环必须执行)
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		target := -1 * nums[i]
		k := n - 1
		// 枚举b
		for j := i + 1; j < n-1; j++ {
			// 需要和上一次枚举的数不同(但是第一次循环必须执行)
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			// 保证b在c左侧
			for j < k && nums[j]+nums[k] > target {
				k--
			}
			// 提前退出
			if j == k {
				break
			}
			// 添加结果
			if nums[j]+nums[k] == target {
				ans = append(ans, []int{nums[i], nums[j], nums[k]})
			}
		}
	}
	return ans
}
