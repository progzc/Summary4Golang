package leetcode_0169_majority_element

import "sort"

// 0169. 多数元素
// https://leetcode.cn/problems/majority-element

// majorityElement 哈希计数
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func majorityElement(nums []int) int {
	n := len(nums)
	maxTimes := 0
	counter := make(map[int]int)
	ans := 0
	for i := 0; i < n; i++ {
		counter[nums[i]]++
		if counter[nums[i]] > maxTimes {
			maxTimes = counter[nums[i]]
			ans = nums[i]
		}
	}
	return ans
}

// majorityElement_2 排序法
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(log(n))
func majorityElement_2(nums []int) int {
	sort.Ints(nums)
	return nums[len(nums)/2]
}

// majorityElement_3 Boyer-Moore 算法（难点是怎么证明这个算法的正确性）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func majorityElement_3(nums []int) int {
	count, candidate := 0, 0
	for _, num := range nums {
		if count == 0 {
			candidate = num
		}
		if num == candidate {
			count += 1
		} else {
			count -= 1
		}
	}
	return candidate
}
