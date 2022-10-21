package leetcode_0485_max_consecutive_ones

// 0485. 最大连续 1 的个数
// https://leetcode.cn/problems/max-consecutive-ones/

// findMaxConsecutiveOnes 一次遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func findMaxConsecutiveOnes(nums []int) int {
	ans, count := 0, 0
	for _, num := range nums {
		if num == 0 {
			ans = max(ans, count)
			count = 0
		} else {
			count++
			ans = max(ans, count)
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
