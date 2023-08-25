package leetcode_0189_rotate_array

// 0189. 轮转数组
// https://leetcode.cn/problems/rotate-array/

// rotate 使用额外的数组
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func rotate(nums []int, k int) {
	newNums := make([]int, len(nums))
	for i, v := range nums {
		newNums[(i+k)%len(nums)] = v
	}
	copy(nums, newNums)
}

// rotate_2 数组翻转
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：先将所有元素翻转，这样尾部的 k%n 个元素就被移至数组头部，然后我们再翻转[0,k%n-1]和[k%n,n-1]区间的元素即能得到最后的答案。
func rotate_2(nums []int, k int) {
	reverse(nums)
	n := len(nums)
	reverse(nums[:k%n])
	reverse(nums[k%n:])
}

func reverse(nums []int) {
	for i, n := 0, len(nums); i < n/2; i++ {
		nums[i], nums[n-1-i] = nums[n-1-i], nums[i]
	}
}
