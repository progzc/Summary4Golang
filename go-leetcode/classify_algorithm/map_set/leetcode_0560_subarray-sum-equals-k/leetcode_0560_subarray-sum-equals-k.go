package leetcode_0560_subarray_sum_equals_k

// 0560.和为K的子数组
// https://leetcode-cn.com/problems/subarray-sum-equals-k/

// subarraySum 枚举
// 时间复杂度: O(n^2)
// 空间复杂度: O(1)
// 思路：注意题目应该是 和为K的连续子数组
//		分别枚举以nums[start]结尾的和为k的连续子数组的个数
func subarraySum(nums []int, k int) int {
	count := 0
	for start := 0; start < len(nums); start++ {
		sum := 0
		for end := start; end >= 0; end-- {
			sum += nums[end]
			if sum == k {
				count++
			}
		}
	}
	return count
}

// subarraySum_2 前缀和+哈希表
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路：注意题目应该是和为K的连续子数组
//		定义pre[i]为[0..i]里所有数的和，则以i结尾的和为k的连续子数组个数时只要统计有多少个前缀和为pre[i]-k的pre[j]即可。
//		构造哈希表：mp[pre[i]−k]int{}，其中键为pre[i]-k，值为出现的次数
func subarraySum_2(nums []int, k int) int {
	count, pre := 0, 0
	m := map[int]int{}
	// 下面这行必须加：不然输入为 [1,1,1] 2 时，会出错
	m[0] = 1 // 含义是第一个元素之前和为0的key对应的值出现的次数默认为1
	for i := 0; i < len(nums); i++ {
		pre += nums[i]
		if _, ok := m[pre-k]; ok {
			count += m[pre-k]
		}
		m[pre] += 1
	}
	return count
}
