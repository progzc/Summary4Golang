package leetcode_1151_minimum_swaps_to_group_all_1s_together

// 1151. 最少交换次数来组合所有的 1
// https://leetcode.cn/problems/minimum-swaps-to-group-all-1s-together/

// minSwaps 滑动窗口
// 思路：假设这里的 1 有 c 个，序列总长度为 n。 我们可以枚举所有长度为 c 的子串，
//		把里面的 0 全部换成 1，所以我们要找到 0 最少的长度为 c 的子串。
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func minSwaps(data []int) int {
	n := len(data)
	totalOne := 0
	for i := 0; i < n; i++ {
		totalOne += data[i]
	}

	countOne := 0
	for i := 0; i < totalOne; i++ {
		countOne += data[i]
	}

	max := countOne
	for i := totalOne; i < n; i++ {
		countOne += data[i] - data[i-totalOne]
		if countOne > max {
			max = countOne
		}
	}
	return totalOne - max
}
