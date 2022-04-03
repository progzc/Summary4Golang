package leetcode_0031_next_permutation

// 0031.下一个排列
// https://leetcode-cn.com/problems/next-permutation/

// nextPermutation 二次遍历
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路:
//	(1)较小数在左边，较大数在右边。
//	(2)较小数尽量靠右，较大数尽量小。
//	(3)交换较大数和较小数后，较大数右边应该升序排列。
func nextPermutation(nums []int) {
	n := len(nums)
	i := n - 2
	// 较小数尽量靠右 且 较小数在左边，较大数在右边
	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	}
	if i >= 0 {
		j := n - 1
		// 较大数尽量小
		for j >= 0 && nums[i] >= nums[j] {
			j--
		}
		// 交换较大数和较小数
		nums[i], nums[j] = nums[j], nums[i]
	}
	// 较大数右边应该升序排列
	//for k := i + 1; k < (n+i+1)/2; k++ {
	//	nums[k], nums[n+i-k] = nums[n+i-k], nums[k]
	//}
	reverse(nums[i+1:])
}

func reverse(nums []int) {
	n := len(nums)
	for i := 0; i < n/2; i++ {
		nums[i], nums[n-1-i] = nums[n-1-i], nums[i]
	}
}
