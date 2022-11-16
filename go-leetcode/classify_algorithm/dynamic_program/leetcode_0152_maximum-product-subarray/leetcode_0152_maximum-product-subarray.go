package leetcode_0152_maximum_product_subarray

import "math"

// 0152. 乘积最大子数组
// https://leetcode.cn/problems/maximum-product-subarray/

// maxProduct
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
// 	a.如果当前位置是一个负数的话，那么我们希望以它前一个位置结尾的某个段的积也是个负数，
//	  这样就可以负负得正，并且我们希望这个积尽可能「负得更多」，即尽可能小。
//	b.如果当前位置是一个正数的话，我们更希望以它前一个位置结尾的某个段的积也是个正数，
//	  并且希望它尽可能地大。
//	状态:
//	  dpmax[i]表示以元素nums[i]结尾的「连续子数组的最大乘积」
//	  dpmin[i]表示以元素nums[i]结尾的「连续子数组的最小乘积」
//	状态转移方程：
//	  当nums[i]<0时:
//	 	 dpmax[i]=max{dpmin[i-1]*nums[i],nums[i]}
//		 dpmin[i]=max{dpmax[i-1]*nums[i],nums[i]}
//	  当nums[i]>=0时: dpmin[i]=max{dpmax[i-1]}
//		 dpmax[i]=max{dpmax[i-1]*nums[i],nums[i]}
//		 dpmin[i]=max{dpmin[i-1]*nums[i],nums[i]}
func maxProduct(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	dpmax, dpmin := make([]int, n), make([]int, n)
	dpmax[0], dpmin[0] = nums[0], nums[0]
	ans := math.MinInt32
	ans = max(ans, dpmax[0])
	for i := 1; i < n; i++ {
		if nums[i] < 0 {
			dpmax[i] = max(dpmin[i-1]*nums[i], nums[i])
			dpmin[i] = min(dpmax[i-1]*nums[i], nums[i])
		} else {
			dpmax[i] = max(dpmax[i-1]*nums[i], nums[i])
			dpmin[i] = min(dpmin[i-1]*nums[i], nums[i])
		}
		ans = max(ans, dpmax[i])
	}
	return ans
}

// maxProduct_2
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路:
// 	a.如果当前位置是一个负数的话，那么我们希望以它前一个位置结尾的某个段的积也是个负数，
//	  这样就可以负负得正，并且我们希望这个积尽可能「负得更多」，即尽可能小。
//	b.如果当前位置是一个正数的话，我们更希望以它前一个位置结尾的某个段的积也是个正数，
//	  并且希望它尽可能地大。
//	状态:
//	  dpmax[i]表示以元素nums[i]结尾的「连续子数组的最大乘积」
//	  dpmin[i]表示以元素nums[i]结尾的「连续子数组的最小乘积」
//	状态转移方程：
//	  当nums[i]<0时:
//	 	 dpmax[i]=max{dpmin[i-1]*nums[i],nums[i]}
//		 dpmin[i]=max{dpmax[i-1]*nums[i],nums[i]}
//	  当nums[i]>=0时: dpmin[i]=max{dpmax[i-1]}
//		 dpmax[i]=max{dpmax[i-1]*nums[i],nums[i]}
//		 dpmin[i]=max{dpmin[i-1]*nums[i],nums[i]}
func maxProduct_2(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	premax, premin := nums[0], nums[0]
	ans := math.MinInt32
	ans = max(ans, premax)
	for i := 1; i < n; i++ {
		var dpmax, dpmin int
		if nums[i] < 0 {
			dpmax = max(premin*nums[i], nums[i])
			dpmin = min(premax*nums[i], nums[i])
		} else {
			dpmax = max(premax*nums[i], nums[i])
			dpmin = min(premin*nums[i], nums[i])
		}
		premax, premin = dpmax, dpmin
		ans = max(ans, dpmax)
	}
	return ans
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
