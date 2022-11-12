package leetcode_0496_next_greater_element_i

// 0496. 下一个更大元素 I
// https://leetcode.cn/problems/next-greater-element-i/

// nextGreaterElement 单调递增(严格)栈
// 时间复杂度: O(m+n)
// 空间复杂度: O(1)
func nextGreaterElement(nums1 []int, nums2 []int) []int {
	m := make(map[int]int, len(nums2))
	var stack []int
	for _, num := range nums2 {
		for len(stack) > 0 && num > stack[len(stack)-1] {
			m[stack[len(stack)-1]] = num
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, num)
	}

	for _, num := range stack {
		m[num] = -1
	}

	var ans []int
	for _, num := range nums1 {
		ans = append(ans, m[num])
	}
	return ans
}
