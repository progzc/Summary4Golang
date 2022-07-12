package leetcode_0370_range_addition

// 370.区间加法
// https://leetcode.cn/problems/range-addition/

// getModifiedArray 常规思路
// 时间复杂度: O(3n)
// 空间复杂度: O(k)
func getModifiedArray(length int, updates [][]int) []int {
	ans := make([]int, length)
	for _, update := range updates {
		for i := update[0]; i <= update[1]; i++ {
			ans[i] += update[2]
		}
	}
	return ans
}

// getModifiedArray_2 差分数组
// 时间复杂度: O(k+n)
// 空间复杂度: O(k)
func getModifiedArray_2(length int, updates [][]int) []int {
	ans := make([]int, length)
	for _, update := range updates {
		ans[update[0]] += update[2]
		if update[1]+1 < length {
			ans[update[1]+1] -= update[2]
		}
	}

	for i := 1; i < length; i++ {
		ans[i] += ans[i-1]
	}
	return ans
}
