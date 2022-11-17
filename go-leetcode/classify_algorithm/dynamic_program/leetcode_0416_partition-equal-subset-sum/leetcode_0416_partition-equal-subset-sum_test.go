package leetcode_0416_partition_equal_subset_sum

import (
	"fmt"
	"testing"
)

// 0416.分割等和子集
// https://leetcode-cn.com/problems/partition-equal-subset-sum/

func Test_canPartition(t *testing.T) {
	nums := []int{1, 5, 11, 15}
	fmt.Println(canPartition(nums))
}

// canPartition 动态规划（0/1背包问题）
// 时间复杂度：O(N*C),其中N为数组元素的个数，C是数组元素的和的一半
// 空间复杂度：O(C)
// 思路：
//	a.数组的和必须为偶数
//	b.可以转换为0/1背包问题：是否可以从输入数组中挑选出一些正整数，使得这些数的和等于整个数组元素的和的一半。
func canPartition(nums []int) bool {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	if sum%2 == 1 {
		return false
	}

	target := sum / 2
	// 状态：dp[i]表示是否存在子集和为i
	dp := make([]bool, target+1)
	// 初始状态：不选择任何元素，则子集和为0
	dp[0] = true
	for _, num := range nums {
		for j := target; j >= num; j-- {
			dp[j] = dp[j] || dp[j-num]
		}
		// 打印中间状态
		fmt.Printf("%v\n", dp)
	}
	return dp[target]
}

// canPartition_2 动态规划（0/1背包问题,二维动态规划）
// 时间复杂度：O(N*C),其中N为数组元素的个数，C是数组元素的和的一半
// 空间复杂度：O(C)
// 思路：
//	a.数组的和必须为偶数
//	b.可以转换为0/1背包问题：是否可以从输入数组中挑选出一些正整数，使得这些数的和等于整个数组元素的和的一半。
// 状态: dp[i][j]表示是否可以从前i个输入数组中挑选出一些正整数，使得这些数的和等于j
func canPartition_2(nums []int) bool {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	if sum%2 == 1 {
		return false
	}

	target := sum / 2
	n := len(nums)
	dp := make([][]bool, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]bool, target+1)
	}

	// 初始状态
	for j := 0; j < target+1; j++ {
		if j == 0 {
			dp[0][j] = true
		} else {
			dp[0][j] = false
		}
	}

	for i := 1; i < n+1; i++ {
		for j := 0; j < target+1; j++ {
			if nums[i-1] > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i-1][j] || dp[i-1][j-nums[i-1]]
			}
		}
	}
	return dp[n][target]
}

//[true	true false false false false false false false false false false false false false false false]
//[true	true false false false true true false false false false false false false false false false]
//[true	true false false false true true false false false false true true false false false true]
//[true	true false false false true true false false false false true true false false true true]
