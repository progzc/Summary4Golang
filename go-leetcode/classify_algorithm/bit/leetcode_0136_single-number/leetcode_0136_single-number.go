package leetcode_0136_single_number

// 0136.只出现一次的数字
// https://leetcode-cn.com/problems/single-number/

// singleNumber 位操作
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func singleNumber(nums []int) int {
	ans := 0
	for _, num := range nums {
		ans ^= num
	}
	return ans
}
