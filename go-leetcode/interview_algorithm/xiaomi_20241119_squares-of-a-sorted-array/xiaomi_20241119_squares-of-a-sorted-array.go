package xiaomi_20241119_squares_of_a_sorted_array

// 小米（一面）
// 有序数组的平方

// sortedSquares 双指针
func sortedSquares(nums []int) []int {
	ans := make([]int, len(nums))
	left, right := 0, len(nums)-1
	i := len(nums) - 1
	for left <= right {
		if nums[left]*nums[left] >= nums[right]*nums[right] {
			ans[i] = nums[left] * nums[left]
			left++
		} else {
			ans[i] = nums[right] * nums[right]
			right--
		}
		i--
	}
	return ans
}
