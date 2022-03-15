package leetcode_0137_single_number_ii

// 0137.只出现一次的数字 II
// https://leetcode-cn.com/problems/single-number-ii/
// 概述：有一个元素出现1次，其他元素都出现3次

// singleNumber 位操作
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 注意事项: 注意限定类型是int32，不然会出错
func singleNumber(nums []int) int {
	ans := int32(0)
	for i := 0; i < 32; i++ {
		total := int32(0)
		for _, num := range nums {
			total += int32(num) >> i & 1
		}
		if total%3 > 0 {
			ans |= 1 << i
		}
	}
	return int(ans)
}
