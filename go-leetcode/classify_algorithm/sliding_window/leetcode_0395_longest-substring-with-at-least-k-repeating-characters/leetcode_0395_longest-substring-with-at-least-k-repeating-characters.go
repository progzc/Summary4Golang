package leetcode_0395_longest_substring_with_at_least_k_repeating_characters

// 0395. 至少有 K 个重复字符的最长子串
// https://leetcode.cn/problems/longest-substring-with-at-least-k-repeating-characters/

func longestSubstring(s string, k int) int {
	// TODO
	return 0
}

func subarraySum(nums []int, k int) int {
	var (
		ans int
		sum int
	)
	for l, r := 0, 0; r < len(nums); r++ {
		sum += nums[r]
		if sum == k {
			ans++
			sum -= nums[l]
			l++
		} else if sum < k {

		}

	}
	return ans
}
