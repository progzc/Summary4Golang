package leetcode_0360_sort_transformed_array

// 0360. 有序转化数组
// https://leetcode.cn/problems/sort-transformed-array/

// sortTransformedArray 抛物线原理+双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：
// 二元一次方程的值，是一个抛物线，但需要考虑 a 和 b 的取值
//	1.若 a > 0, 最小点在 mid=-1*b/2/a, x 离 mid 越远，函数值越大
//	2.如 a < 0, 最大点在 mid=-1*b/2/a, x 离 mid 越远，函数值越小
//	3.若 a == 0 && b >= 0, 函数值单调递增
//	4.若 a == 0 && b < 0, 函数值单调递减
func sortTransformedArray(nums []int, a int, b int, c int) []int {
	n := len(nums)
	if n == 0 {
		return nil
	}

	ans := make([]int, n)
	if a == 0 {
		// 3.若 a == 0 && b >= 0, 函数值单调递增
		for i := 0; i < n; i++ {
			ans[i] = calc(nums[i], a, b, c)
		}
		// 4.若 a == 0 && b < 0, 函数值单调递减
		if b < 0 {
			reverse(ans)
		}
		return ans
	}

	mid, idx := -1*float32(b)/2/float32(a), n-1
	left, right := 0, len(nums)-1
	for left <= right {
		// 1.若 a > 0, 最小点在 mid=-1*b/2/a, x 离 mid 越远，函数值越大
		if abs(mid-float32(nums[left])) >= abs(float32(nums[right])-mid) {
			ans[idx] = calc(nums[left], a, b, c)
			left++
		} else {
			ans[idx] = calc(nums[right], a, b, c)
			right--
		}
		idx--
	}

	// 2.如 a < 0, 最大点在 mid=-1*b/2/a, x 离 mid 越远，函数值越小
	if a < 0 {
		reverse(ans)
	}
	return ans
}

func abs(x float32) float32 {
	if x > 0 {
		return x
	}
	return -x
}

func calc(x, a, b, c int) int {
	return a*x*x + b*x + c
}

func reverse(nums []int) {
	n := len(nums)
	for i := 0; i < n/2; i++ {
		nums[i], nums[n-1-i] = nums[n-1-i], nums[i]
	}
}
