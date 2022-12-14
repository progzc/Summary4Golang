package leetcode_0698_partition_to_k_equal_sum_subsets

import "sort"

// 698. 划分为k个相等的子集
// https://leetcode.cn/problems/partition-to-k-equal-sum-subsets/

// canPartitionKSubsets dfs (超时)
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(2^n)
func canPartitionKSubsets(nums []int, k int) bool {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	if sum%k != 0 {
		return false
	}
	target := sum / k
	bucket := make([]int, k)
	n := len(nums)

	var (
		dfs func(idx int) bool
	)
	dfs = func(idx int) bool {
		if idx == n {
			for i := 0; i < k; i++ {
				if bucket[i] != target {
					return false
				}
			}
			return true
		}

		for i := 0; i < k; i++ {
			// 剪枝：放入球后超过 target 的值，选择一下桶
			if bucket[i]+nums[idx] > target {
				continue
			}
			// 做选择：放入 i 号桶
			bucket[i] += nums[idx]
			// 处理下一个数
			if dfs(idx + 1) {
				return true
			}
			// 撤销选择：挪出i号桶
			bucket[i] -= nums[idx]
		}
		return false
	}
	return dfs(0)
}

// canPartitionKSubsets_2 dfs (优化)
// 时间复杂度: O(n*2^n)
// 空间复杂度: O(2^n)
func canPartitionKSubsets_2(nums []int, k int) bool {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	if sum%k != 0 {
		return false
	}
	target := sum / k
	bucket := make([]int, k)
	n := len(nums)

	// 降序排列：先让值大的元素选择桶，这样可以增加剪枝的命中率，从而降低回溯的概率
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] > nums[j]
	})

	var (
		dfs func(idx int) bool
	)
	dfs = func(idx int) bool {
		if idx == n {
			// 其实这个地方不需要判断，因为当 idx == len(nums) 时，所有球已经按要求装入所有桶，所以肯定是一个满足要求的解
			// 即：每个桶内球的和一定为 target。
			//for i := 0; i < k; i++ {
			//	if bucket[i] != target {
			//		return false
			//	}
			//}
			return true
		}

		for i := 0; i < k; i++ {
			// 如果当前桶和上一个桶内的元素和相等，则跳过。
			// 原因：如果元素和相等，那么 nums[index] 选择上一个桶和选择当前桶可以得到的结果是一致的
			if i > 0 && bucket[i] == bucket[i-1] {
				continue
			}
			// 剪枝：放入球后超过 target 的值，选择一下桶
			if bucket[i]+nums[idx] > target {
				continue
			}
			// 做选择：放入 i 号桶
			bucket[i] += nums[idx]
			// 处理下一个数
			if dfs(idx + 1) {
				return true
			}
			// 撤销选择：挪出i号桶
			bucket[i] -= nums[idx]
		}
		return false
	}
	return dfs(0)
}
