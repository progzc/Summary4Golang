package leetcode_0992_subarrays_with_k_different_integers

// 0992. K 个不同整数的子数组
// https://leetcode.cn/problems/subarrays-with-k-different-integers/

// subarraysWithKDistinct 普通解法 (超时)
// 时间复杂度: O(n^2)
// 空间复杂度: O(n)
func subarraysWithKDistinct(nums []int, k int) int {
	ans := 0
	n := len(nums)
	if n < k || k == 0 {
		return ans
	}

	for l := 0; l <= n-k; l++ {
		m := make(map[int]int, 0)
		r := l
		for r < n && len(m) <= k {
			m[nums[r]]++
			if len(m) == k {
				ans++
			}
			r++
		}
	}
	return ans
}

// subarraysWithKDistinct_2 滑动窗口
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路: 我们可以证明：对于任意一个右端点，能够与其对应的左端点们必然相邻。
func subarraysWithKDistinct_2(nums []int, k int) int {
	// TODO
	return 0
}
