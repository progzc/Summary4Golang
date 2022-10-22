package leetcode_1493_longest_subarray_of_1s_after_deleting_one_element

// 1493. 删掉一个元素以后全为 1 的最长子数组
// https://leetcode.cn/problems/longest-subarray-of-1s-after-deleting-one-element/

// longestSubarray 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func longestSubarray(nums []int) int {
	n := len(nums)
	ans := 0
	freq := [2]int{}
	for l, r := 0, 0; r < n; r++ {
		freq[nums[r]]++
		for freq[0] > 1 {
			freq[nums[l]]--
			l++
		}
		ans = max(ans, r-l)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
