package leetcode_0239_sliding_window_maximum

// 0239.滑动窗口最大值
// https://leetcode-cn.com/problems/sliding-window-maximum/

// maxSlidingWindow 单调递减队列
// 时间复杂度: O(n)
// 空间复杂度: O(k)
// 思路:
//	我们可以使用一个队列存储所有还没有被移除的下标。在队列中，这些下标按照从小到大的顺序被存储，并且它们在数组 nums 中对应的值是严格单调递减的。
//	简言之, 队列中随着下标递增, 值递减。
func maxSlidingWindow(nums []int, k int) []int {
	var q []int
	push := func(i int) {
		for len(q) > 0 && nums[q[len(q)-1]] <= nums[i] {
			q = q[:len(q)-1]
		}
		q = append(q, i)
	}
	for i := 0; i < k; i++ {
		push(i)
	}
	n := len(nums)
	ans := make([]int, 1, n-k+1)
	ans[0] = nums[q[0]]
	for i := k; i < n; i++ {
		push(i)
		for q[0] <= i-k {
			q = q[1:]
		}
		ans = append(ans, nums[q[0]])
	}
	return ans
}
