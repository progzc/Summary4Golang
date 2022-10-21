package leetcode_1004_max_consecutive_ones_iii

// 1004. 最大连续1的个数 III
// https://leetcode.cn/problems/max-consecutive-ones-iii/

// 与 https://leetcode.cn/problems/max-consecutive-ones-ii/ 类似

// longestOnes 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：进行题意转换，把「最多可以把 K 个 0 变成 1，求仅包含 1 的最长子数组的长度」
//		转换为 「找出一个最长的子数组，该子数组内最多允许有 K 个 0 」
func longestOnes(nums []int, k int) int {
	n := len(nums)
	count, ans := 0, 0
	for l, r := 0, 0; r < n; r++ {
		if nums[r] == 0 {
			count++
		}
		for count > k {
			if nums[l] == 0 {
				count--
				l++
			} else {
				l++
			}
		}
		ans = max(ans, r-l+1)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
