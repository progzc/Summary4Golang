package leetcode_0128_longest_consecutive_sequence

import "sort"

// 0128. 最长连续序列
// https://leetcode.cn/problems/longest-consecutive-sequence/

// longestConsecutive 排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
// 注意事项：注意nums中可能存在相等的数字
//	例：
//		输入[1,2,0,1],输出3
//		输入[],输出0
func longestConsecutive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	ans, count := 1, 1
	for i := 0; i < len(nums); i++ {
		if i > 0 {
			if nums[i] == nums[i-1] {
				continue
			} else if nums[i] == nums[i-1]+1 {
				count++
				if count > ans {
					ans = count
				}
			} else {
				count = 1
			}
		}
	}
	return ans
}

// longestConsecutive_2 哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func longestConsecutive_2(nums []int) int {
	set := map[int]bool{}
	for _, num := range nums {
		set[num] = true
	}

	ans := 0
	for num := range set {
		// 这个if很关键,直接降低了时间复杂度
		if !set[num-1] {
			cur := num
			count := 1
			for set[cur+1] {
				count++
				cur++
			}
			if count > ans {
				ans = count
			}
		}
	}
	return ans
}
