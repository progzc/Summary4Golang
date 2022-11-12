package leetcode_0503_next_greater_element_ii

// 0503. 下一个更大元素 II
// https://leetcode.cn/problems/next-greater-element-ii/

// nextGreaterElements 单调递减(严格)栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：两次搜索
func nextGreaterElements(nums []int) []int {
	n := len(nums)
	var (
		ans   = make([]int, 2*n)
		stack []int
	)
	nums = append(nums, nums...)
	for i, num := range nums {
		for len(stack) > 0 && num > nums[stack[len(stack)-1]] {
			ans[stack[len(stack)-1]] = num
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	for _, v := range stack {
		ans[v] = -1
	}
	return ans[0:n]
}
