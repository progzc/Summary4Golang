package leetcode_0018_4sum

import "sort"

// 18. 四数之和
// https://leetcode.cn/problems/4sum/

// fourSum 排序+双指针
// 时间复杂度: O(n^3)
// 空间复杂度: O(log(n))
// 思路：
// 1. 凡是不重复，首先要想到先排序 或 采用散列表
// 2. 排序之后采用四重循环，要保证每重循环 相邻两次枚举的元素不能相同
// 双指针：当随着第一个元素的递增,第二个元素是递减的，就可以考虑使用 双指针
func fourSum(nums []int, target int) [][]int {
	var ans [][]int
	n := len(nums)
	if n < 4 {
		return ans
	}

	sort.Ints(nums)
	for i := 0; i < n-3; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		for j := i + 1; j < n-2; j++ {
			if j > i+1 && nums[j] == nums[j-1] {
				continue
			}
			s := n - 1
			newTarget := target - (nums[i] + nums[j])
			for k := j + 1; k < n-1; k++ {
				if k > j+1 && nums[k] == nums[k-1] {
					continue
				}
				for k < s && nums[k]+nums[s] > newTarget {
					s--
				}
				if k == s {
					break
				}
				if nums[k]+nums[s] == newTarget {
					ans = append(ans, []int{nums[i], nums[j], nums[k], nums[s]})
				}
			}
		}
	}
	return ans
}
