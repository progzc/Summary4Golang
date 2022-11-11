package leetcode_0042_trapping_rain_water

// 0042.接雨水
// https://leetcode-cn.com/problems/trapping-rain-water/

// trap 动态规划
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：每个位置接的雨水=左右最大高度的最小值-当前高度
func trap(height []int) int {
	ans := 0
	n := len(height)
	if n == 0 {
		return ans
	}

	leftMax, rightMax := make([]int, n), make([]int, n)
	// 这里是关键
	leftMax[0], rightMax[n-1] = height[0], height[n-1]

	for i := 1; i < n; i++ {
		leftMax[i] = max(leftMax[i-1], height[i])
	}

	for i := n - 2; i >= 0; i-- {
		rightMax[i] = max(rightMax[i+1], height[i])
	}

	for i := 0; i < n; i++ {
		ans += min(leftMax[i], rightMax[i]) - height[i]
	}
	return ans
}

// trap_2 单调递减(非严格)栈
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路：画出柱状图+想到用单调递减(非严格)栈
func trap_2(height []int) int {
	ans := 0
	stack := []int{}
	for i, h := range height {
		for len(stack) > 0 && h > height[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				break
			}
			left := stack[len(stack)-1]
			curW := i - left - 1
			curH := min(height[left], h) - height[top]
			ans += curW * curH
		}
		// 保证每个都有入栈的机会
		stack = append(stack, i)
	}
	return ans
}

// trap_3 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：这里有点难以理解
func trap_3(height []int) int {
	ans := 0
	left, right := 0, len(height)-1
	leftMax, rightMax := 0, 0
	for left < right {
		leftMax = max(leftMax, height[left])
		rightMax = max(rightMax, height[right])
		if height[left] < height[right] {
			ans += leftMax - height[left]
			left++
		} else {
			ans += rightMax - height[right]
			right--
		}
	}
	return ans
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
