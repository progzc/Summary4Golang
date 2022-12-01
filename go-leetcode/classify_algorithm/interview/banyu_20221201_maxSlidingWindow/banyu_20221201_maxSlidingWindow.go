package banyu_20221201_maxSlidingWindow

// maxSlidingWindow 单调队列
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func maxSlidingWindow(nums []int, k int) []int {
	var ans []int
	n := len(nums)
	if n < k {
		return ans
	}

	var (
		push  func(i int)
		stack []int
	)

	push = func(i int) {
		for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[i] {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	for i := 0; i < k; i++ {
		push(i)
	}

	ans = append(ans, nums[stack[0]])
	for i := k; i < n; i++ {
		push(i)
		for stack[0] <= i-k {
			stack = stack[1:]
		}
		ans = append(ans, nums[stack[0]])
	}
	return ans
}
