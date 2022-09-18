package leetcode_1060_missing_element_in_sorted_array

// 1060. 有序数组中的缺失元素
// https://leetcode.cn/problems/missing-element-in-sorted-array/

// missingElement 顺序法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func missingElement(nums []int, k int) int {
	omit := 0
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1]+1 {
			omit += nums[i] - nums[i-1] - 1
			if omit >= k {
				return nums[i] - 1 - (omit - k)
			}
		}
	}
	return nums[len(nums)-1] + (k - omit)
}

// missingElement_2 二分法
// 时间复杂度: O(log(n))
// 空间复杂度: O(1)
// 思路：
//	对于输入nums=[4, 7, 9, 10],k=3,miss数组=[0 ,2 , 3, 3]
//	问题转化为在miss数组中找到第一个大于等于k的下标idx
//	则缺失的元素是 nums[idx] - nums[0] - idx
func missingElement_2(nums []int, k int) int {
	n := len(nums)
	if k > miss(n-1, nums) {
		return nums[n-1] + k - miss(n-1, nums)
	}

	left, right := 0, n-1
	pos := -1
	for left <= right {
		mid := left + (right-left)/2
		if miss(mid, nums) >= k {
			// 找到第一个大于等于 target 的元素
			if mid == 0 || miss(mid-1, nums) < k {
				pos = mid
				break
			}
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	// 根据题意：miss数组长度大于0,miss[0]=0,且k>0；
	// 这里必然有的pos>=1，不会出现数组越界情况
	return nums[pos-1] + k - miss(pos-1, nums)
}

func miss(idx int, nums []int) int {
	return nums[idx] - nums[0] - idx
}
