package leetcode_0581_shortest_unsorted_continuous_subarray

import (
	"sort"
)

// 0581. 最短无序连续子数组
// https://leetcode.cn/problems/shortest-unsorted-continuous-subarray/

// findUnsortedSubarray 排序
// 时间复杂度: O(nlog(n))
// 空间复杂度: O(n)
// 思路:
// 	我们将给定的数组 nums 表示为三段子数组拼接的形式,分别记作 numsA, numsB, numsC。当我们对 numsB
// 	进行排序，整个数组将变为有序。换而言之，当我们对整个序列进行排序，numsA 和 numsC 都不会改变。
func findUnsortedSubarray(nums []int) int {
	if sort.IntsAreSorted(nums) {
		return 0
	}
	numsSorted := append([]int(nil), nums...)
	sort.Ints(numsSorted)
	l, r := 0, len(numsSorted)-1
	for nums[l] == numsSorted[l] {
		l++
	}
	for nums[r] == numsSorted[r] {
		r--
	}
	return r - l + 1
}

// findUnsortedSubarray_2 单调栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	i)利用单调栈求出排序区间（题目中的连续子数组）的左右边界。注意：栈中存放的是下标。
//	ii)首先求左边界，从前往后遍历数组，维护一个单调递增栈。当遍历到的元素 nums[i] 比栈顶元素 nums[stack[stack.length-1]] 小时，
//	   不断将栈顶元素踢出，直到当前元素比栈顶元素大或者相等。 那么这样做有什么用呢？仔细思考一下，对于栈中的元素，
//	   它的下标一定比当前遍历到的元素的下标小，而当前元素比栈顶元素小时，说明排序之后栈顶元素的下标肯定是会变化的，
//	   也就是栈顶元素一定会出现在题目要求我们找的那个连续子数组中。 这个时候我们只需要记录当前被踢出栈顶元素的下标，找出最左边的那个被踢出栈的元素，
//	   它就是题目中的连续子数组的左边界。
//	iii)然后求右边界，从后往前遍历数组，维护一个单调递减栈。这个与求左边界是同样的道理，只是栈变为了单调递减栈，因为我们要使得数组变成升序序列。
func findUnsortedSubarray_2(nums []int) int {
	var (
		s1 []int
		// 注意：若这里l = len(nums)-1,则针对如下用例会报错
		// 输入: nums=[1]; 输出: 0
		l = len(nums)
	)
	for i, num := range nums {
		for len(s1) > 0 && num < nums[s1[len(s1)-1]] {
			l = min(l, s1[len(s1)-1])
			s1 = s1[:len(s1)-1]
		}
		s1 = append(s1, i)
	}

	var (
		s2 []int
		r  = -1
	)
	for i := len(nums) - 1; i >= 0; i-- {
		for len(s2) > 0 && nums[i] > nums[s2[len(s2)-1]] {
			r = max(r, s2[len(s2)-1])
			s2 = s2[:len(s2)-1]
		}
		s2 = append(s2, i)
	}

	if r-l+1 > 0 {
		return r - l + 1
	}
	return 0
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
