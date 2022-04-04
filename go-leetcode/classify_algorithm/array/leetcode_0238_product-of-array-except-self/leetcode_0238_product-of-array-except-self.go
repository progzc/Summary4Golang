package leetcode_0238_product_of_array_except_self

// 238.除自身以外数组的乘积
// https://leetcode-cn.com/problems/product-of-array-except-self/

// productExceptSelf 左右乘积列表
// 时间复杂度：O(n)
// 空间复杂度：O(n)
// 思路：以空间换时间
//	利用索引左侧所有数字的乘积和右侧所有数字的乘积（即前缀与后缀）相乘。
//	L[i]：代表的是 i 左侧（不包括nums[i]）所有数字的乘积。
//	R[i]：代表的是 i 右侧（不包括nums[i]）所有数字的乘积。
func productExceptSelf(nums []int) []int {
	length := len(nums)
	L, R, ans := make([]int, length), make([]int, length), make([]int, length)
	L[0], R[length-1] = 1, 1

	for i := 1; i < length; i++ {
		L[i] = nums[i-1] * L[i-1]
	}

	for i := length - 2; i >= 0; i-- {
		R[i] = nums[i+1] * R[i+1]
	}

	for i := 0; i < length; i++ {
		ans[i] = L[i] * R[i]
	}
	return ans
}

// productExceptSelf_2 优化空间复杂度
// 时间复杂度：O(n)
// 空间复杂度：O(1)
// 思路：利用输出数组作为存储
func productExceptSelf_2(nums []int) []int {
	length := len(nums)
	ans := make([]int, length)

	ans[0] = 1
	for i := 1; i < length; i++ {
		ans[i] = nums[i-1] * ans[i-1]
	}

	R := 1
	for i := length - 1; i >= 0; i-- {
		ans[i] = ans[i] * R
		R = R * nums[i]
	}
	return ans
}
