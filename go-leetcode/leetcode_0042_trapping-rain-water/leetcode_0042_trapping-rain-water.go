package leetcode_0042_trapping_rain_water

import "github.com/progzc/Summary4Golang/go-leetcode/structures"

// 42. 接雨水
// link: https://leetcode-cn.com/problems/trapping-rain-water/

// trap1 暴力法
// 时间复杂度 O(n^2)
// 空间复杂度 O(1)
func trap1(height []int) int {
	ans, size := 0, len(height)
	for i := 1; i < size-1; i++ {
		maxLeft, maxRight := 0, 0
		for j := i; j >= 0; j-- {
			maxLeft = max(maxLeft, height[j])
		}
		for j := i; j < size; j++ {
			maxRight = max(maxRight, height[j])
		}
		ans += min(maxLeft, maxRight) - height[i]
	}
	return ans
}

// trap2 动态编程（对暴力法的优化）
// 时间复杂度 O(n)
// 空间复杂度 O(n)
func trap2(height []int) int {
	if height == nil || len(height) == 0 {
		return 0
	}
	ans, size := 0, len(height)
	maxLeft, maxRight := make([]int, size), make([]int, size)

	// 1. 找到数组中从下标 i 到最左端最高的条形块高度maxLeft
	maxLeft[0] = height[0]
	for i := 1; i < size; i++ {
		maxLeft[i] = max(height[i], maxLeft[i-1])
	}
	// 2. 找到数组中从下标 i 到最右端最高的条形块高度maxRight
	maxRight[size-1] = height[size-1]
	for i := size - 2; i >= 0; i-- {
		maxRight[i] = max(height[i], maxRight[i+1])
	}
	// 3. 扫描数组height并更新答案
	for i := 1; i < size-1; i++ {
		ans += min(maxLeft[i], maxRight[i]) - height[i]
	}
	return ans
}

// trap3 单调栈
// 时间复杂度 O(n)
// 空间复杂度 O(n)
func trap3(height []int) int {
	ans, current := 0, 0
	stack := structures.NewStack()
	for current < len(height) {
		for !stack.IsEmpty() && height[current] > height[stack.Peek()] {
			top := stack.Pop()
			if stack.IsEmpty() {
				break
			}
			distance := current - stack.Peek() - 1
			heightDiff := min(height[current], height[stack.Peek()]) - height[top]
			ans += distance * heightDiff
		}
		stack.Push(current)
		current++
	}
	return ans
}

// trap4 双指针
// 时间复杂度 O(n)
// 空间复杂度 O(1)
func trap4(height []int) int {
	left, right, ans := 0, len(height)-1, 0
	maxLeft, maxRight := 0, 0
	for left < right {
		if height[left] < height[right] {
			if height[left] >= maxLeft {
				maxLeft = height[left]
			} else {
				ans += maxLeft - height[left]
			}
			left++
		} else {
			if height[right] >= maxRight {
				maxRight = height[right]
			} else {
				ans += maxRight - height[right]
			}
			right--
		}
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
