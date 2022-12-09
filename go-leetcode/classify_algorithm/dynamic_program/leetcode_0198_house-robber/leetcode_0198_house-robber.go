package leetcode_0198_house_robber

// 198. 打家劫舍
// https://leetcode.cn/problems/house-robber/

// rob 动态规划（二维）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	状态: dp[i][j]表示为从前i个房间能够偷盗的最高金额;
//		 其中j=0表示不偷盗第i个房间, j=1表示偷盗第i个房间
//	转移方程:
//		dp[i+1][0] = max(dp[i][0],dp[i][1])
//		dp[i+1][1] = dp[i][0]+nums[i+1]
func rob(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, 2)
	}
	dp[0][0], dp[0][1] = 0, nums[0]
	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][0], dp[i-1][1])
		dp[i][1] = dp[i-1][0] + nums[i]
	}
	return max(dp[n-1][0], dp[n-1][1])
}

// rob_2 动态规划（常量）
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 思路:
//	状态: dp[i][j]表示为从前i个房间能够偷盗的最高金额;
//		 其中j=0表示不偷盗第i个房间, j=1表示偷盗第i个房间
//	转移方程:
//		dp[i+1][0] = max(dp[i][0],dp[i][1])
//		dp[i+1][1] = dp[i][0]+nums[i+1]
func rob_2(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nums[0]
	}

	dp0, dp1 := 0, nums[0]
	for i := 1; i < n; i++ {
		newDp0 := max(dp0, dp1)
		newDp1 := dp0 + nums[i]
		dp0, dp1 = newDp0, newDp1
	}
	return max(dp0, dp1)
}

// rob_3 动态规划（一维）
// 时间复杂度: O(n)
// 空间复杂度: O(n)
// 思路:
//	状态: dp[i]表示为从前i个房间能够偷盗的最高金额;
//	转移方程:
//		dp[i] = max(dp[i-2]+nums[i], dp[i-1])
//	边界条件:
//		dp[0] = nums[0] 				// 只有一间房屋，则偷窃该房屋
//		dp[1] = max(nums[0],nums[1]) 	// 只有两间房屋，选择其中金额较高的房屋进行偷窃
func rob_3(nums []int) int {
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
