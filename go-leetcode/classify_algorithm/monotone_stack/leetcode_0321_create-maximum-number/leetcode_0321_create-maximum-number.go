package leetcode_0321_create_maximum_number

// 0321. 拼接最大数
// https://leetcode.cn/problems/create-maximum-number/

// 类似题:
// 0402. 移掉 K 位数字
// https://leetcode.cn/problems/remove-k-digits/

// maxNumber 单调递减(非严格)栈+合并
// 时间复杂度: O(k*(m+n+k^2))
// 空间复杂度: O(n)
// 思路:
//	为了找到长度为 k 的最大数，需要从两个数组中分别选出最大的子序列，这两个子序列的长度之和为 k，然后将这两个子序列合并得到最大数。
//	两个子序列的长度最小为 0，最大不能超过 k 且不能超过对应的数组长度。
//		i)第一步可以通过单调栈实现。
//		ii)第二步可以通过合并实现
func maxNumber(nums1 []int, nums2 []int, k int) []int {
	start := 0
	if k > len(nums2) {
		start = k - len(nums2)
	}
	var ans []int
	for i := start; i <= k && i <= len(nums1); i++ {
		s1 := maxSubsequence(nums1, i)
		s2 := maxSubsequence(nums2, k-i)
		merged := merge(s1, s2)
		if compareLess(ans, merged) {
			ans = merged
		}
	}
	return ans
}

// maxSubsequence 单调递减(非严格)栈
// 功能: 求nums中长度为k的最大子序列
// 思路:
// 	0402. 移掉 K 位数字
//	https://leetcode.cn/problems/remove-k-digits/
func maxSubsequence(nums []int, k int) []int {
	n := len(nums)
	if n <= k {
		return nums
	}

	var (
		stack []int
		del   = n - k
	)
	for _, digit := range nums {
		for del > 0 && len(stack) > 0 && digit > stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
			del--
		}
		stack = append(stack, digit)
	}
	stack = stack[:len(stack)-del]
	return stack
}

func merge(nums1, nums2 []int) []int {
	merged := make([]int, len(nums1)+len(nums2))
	for i := range merged {
		if compareLess(nums1, nums2) {
			merged[i], nums2 = nums2[0], nums2[1:]
		} else {
			merged[i], nums1 = nums1[0], nums1[1:]
		}
	}
	return merged
}

func compareLess(nums1, nums2 []int) bool {
	for i := 0; i < len(nums1) && i < len(nums2); i++ {
		if nums1[i] != nums2[i] {
			return nums1[i] < nums2[i]
		}
	}
	return len(nums1) < len(nums2)
}
