package leetcode_0213_house_robber_ii

// 213. 打家劫舍 II
// https://leetcode.cn/problems/house-robber-ii/

// rob
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	假设数组 nums 的长度为 n。如果不偷窃最后一间房屋，则偷窃房屋的下标范围是 [0, n-2];
//	如果不偷窃第一间房屋，则偷窃房屋的下标范围是 [1, n-1]。
//	则只需要两次执行【198. 打家劫舍】即可。
func rob(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	return max(_rob(nums[0:n-1]), _rob(nums[1:n]))
}

// _rob 动态规划（一维）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 注意: 同下面解法
//	198. 打家劫舍
//	https://leetcode.cn/problems/house-robber/
// 思路:
//	状态: dp[i]表示为从前i个房间能够偷盗的最高金额;
//	转移方程:
//		dp[i] = max(dp[i-2]+nums[i], dp[i-1])
//	边界条件:
//		dp[0] = nums[0] 				// 只有一间房屋，则偷窃该房屋
//		dp[1] = max(nums[0],nums[1]) 	// 只有两间房屋，选择其中金额较高的房屋进行偷窃
func _rob(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	dp := make([]int, n)
	dp[0] = nums[0]
	dp[1] = max(nums[0], nums[1])
	for i := 2; i < n; i++ {
		dp[i] = max(dp[i-2]+nums[i], dp[i-1])
	}
	return dp[n-1]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
