package leetcode_0011_container_with_most_water

// 0011.盛最多水的容器
// https://leetcode-cn.com/problems/container-with-most-water/

// maxArea 双指针
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：贪心
func maxArea(height []int) int {
	left, right, area := 0, len(height)-1, 0
	for left < right {
		if height[left] <= height[right] {
			area = max(area, height[left]*(right-left))
			left++
		} else {
			area = max(area, height[right]*(right-left))
			right--
		}
	}
	return area
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
