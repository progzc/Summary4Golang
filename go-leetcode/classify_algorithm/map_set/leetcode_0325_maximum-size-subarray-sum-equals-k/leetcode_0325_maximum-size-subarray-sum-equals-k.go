package leetcode_0325_maximum_size_subarray_sum_equals_k

// 0325. 和等于 k 的最长子数组长度
// https://leetcode.cn/problems/maximum-size-subarray-sum-equals-k/

// maxSubArrayLen 前缀和+哈希
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func maxSubArrayLen(nums []int, k int) int {
	n := len(nums)
	m := make(map[int]int)

	ans, preSum := 0, 0
	// 注意：0出现在位置为-1位置处
	m[0] = -1
	for i := 0; i < n; i++ {
		// 累加前缀和
		preSum += nums[i]
		// 确保出现的是第一次出现的位置
		if _, ok := m[preSum]; !ok {
			m[preSum] = i
		}
		// 每次检查下是否需要更新答案
		if v, ok := m[preSum-k]; ok {
			ans = max(ans, i-v)
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
